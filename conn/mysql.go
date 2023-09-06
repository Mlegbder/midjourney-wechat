package conn

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"mj-wechat-bot/config"
	"mj-wechat-bot/errorhandler"
	"time"
)

var DB *gorm.DB

func init() {
	// 注册异常处理函数
	defer errorhandler.HandlePanic()
	// Read configuration file.
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// Unmarshal configuration.
	var msl config.Config
	err = yaml.Unmarshal(data, &msl)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	DB, err = gorm.Open(mysql.Open(msl.Mysql.Dsn), &gorm.Config{
		SkipDefaultTransaction: true, //初始化时 true 禁用它，这将获得大约 30%+ 性能提升。
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("DB connect failed:" + err.Error())
	}
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(msl.Mysql.MaxIdleConn)
	db.SetMaxOpenConns(msl.Mysql.MaxOpenConn)
	db.SetConnMaxLifetime(time.Hour)
	fmt.Println("Mysql Connect Successful！")
}

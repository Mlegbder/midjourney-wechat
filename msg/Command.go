package msg

import (
	"encoding/base64"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"log"
	"mj-wechat-bot/api"
	"mj-wechat-bot/conn"
	"mj-wechat-bot/consts"
	"mj-wechat-bot/model"
	"mj-wechat-bot/replay"
	"mj-wechat-bot/task"
	"mj-wechat-bot/utils"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	Commands = map[string]string{
		"/imagine":  "Imagine",
		"/up":       "Up",
		"/help":     "Help",
		"/cz":       "Recharge",
		"/me":       "UserInfo",
		"/describe": "Describe",
		"/blend":    "Blend",
		"/sucai":    "Parse",
	}
)

type Command interface {
	Imagine()
	Up()
	Help()
	UserInfo()
	Describe()
	Blend()
	Parse()
}
type Impl struct {
	msg     *openwechat.Message
	realMsg string
	info    replay.Info
}

func (c Impl) call(pre string, command string) {
	c.realMsg = strings.ReplaceAll(c.realMsg, pre, "")
	c.realMsg = strings.TrimSpace(c.realMsg)
	log.Printf("调用命令: %s,内容: %s\n", command, c.realMsg)
	// 获取结构体反射对象
	function := reflect.ValueOf(c)
	//log.Printf("impl:%v", function)
	// 获取结构体方法的反射对象
	method := function.MethodByName(command)
	//log.Printf("method:%v", method)
	// 调用方法
	method.Call(nil)
}

func (c Impl) Imagine() {
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	if c.realMsg == "" {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskMainCommandErrMsg))
		return
	}
	if c.CheckUserBalance(consts.IMAGINE) {
		ok, taskId := api.CreateMessage(c.realMsg)
		if ok {
			c.info.TaskId = taskId
			c.info.Prompt = c.realMsg
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskMainCreateMsg))
			log.Printf("任务已经提交:%s", taskId)
			c.msg.Set("type", "main")
			task.AddTask(c.msg, taskId)
			err = c.DeductionsBalance(consts.IMAGINE)
			if err != nil {
				log.Println(err)
			}
		} else {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSendErrMsg))
		}
	}
}

func (c Impl) Up() {
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	commands := strings.SplitN(c.realMsg, " ", 2)
	if len(commands) != 2 {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSubCommandErrMsg))
		return
	}
	taskId := strings.TrimSpace(commands[0])
	action := strings.ToLower(strings.TrimSpace(commands[1]))
	c.info.TaskId = taskId
	c.info.Action = action
	//判断action是否在指定字符串内
	switch action {
	case "u1", "u2", "u3", "u4", "v1", "v2", "v3", "v4":
		break
	default:
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSubCommandErrMsg))
		//c.msg.ReplyText("参数错误,可选参数:u1,u2,u3,u4,v1,v2,v3,v4")
		return
	}

	flag := true
	if strings.Contains(action, "u") {
		flag = c.CheckUserBalance(consts.UPSCALE)
	} else if strings.Contains(action, "v") {
		flag = c.CheckUserBalance(consts.VARIATION)
	}
	if flag {
		ok, newTaskId := api.TaskUpdate(taskId, action)
		if ok {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSubCreateMsg))
			log.Printf("更新任务已经提交:%s", newTaskId)
			c.msg.Set("type", strings.ToUpper(action))
			task.AddTask(c.msg, newTaskId)
			if strings.Contains(action, "u") {
				err = c.DeductionsBalance(consts.UPSCALE)
				if err != nil {
					log.Println(err)
				}
			} else if strings.Contains(action, "v") {
				err = c.DeductionsBalance(consts.VARIATION)
				if err != nil {
					log.Println(err)
				}
			}

		} else {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSendErrMsg))
			//c.msg.ReplyText("任务创建失败")
		}
	}

}

func (c Impl) Blend() {
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	bsArray := c.GetImgs()
	if err != nil {
		log.Println(err)
		return
	}
	if c.CheckUserBalance(consts.BLEND) {
		ok, taskId := api.Blend(bsArray)
		if ok {
			c.info.TaskId = taskId
			c.info.Prompt = c.realMsg
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskBlendMsg))
			log.Printf("任务已经提交:%s", taskId)
			c.msg.Set("type", "BLEND")
			task.AddTask(c.msg, taskId)
			err = c.DeductionsBalance(consts.BLEND)
			if err != nil {
				log.Println(err)
			}
		} else {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSendErrMsg))
		}
	}
}

func (c Impl) Describe() {
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	base64, err := c.GetImg()
	if err != nil {
		log.Println(err)
		return
	}
	if c.CheckUserBalance(consts.DESCRIBE) {
		ok, taskId := api.Describe(base64)
		if ok {
			c.info.TaskId = taskId
			c.info.Prompt = c.realMsg
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskDescribeMsg))
			log.Printf("任务已经提交:%s", taskId)
			c.msg.Set("type", "DESCRIBE")
			task.AddTask(c.msg, taskId)
			err = c.DeductionsBalance(consts.DESCRIBE)
			if err != nil {
				log.Println(err)
			}
		} else {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSendErrMsg))
		}
	}
}

func (c Impl) Parse() {
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	if c.CheckUserBalance(consts.PARSE) {
		ok, resp := api.Parse(c.realMsg)
		if ok {
			c.info.Url = resp
			err = c.DeductionsBalance(consts.PARSE)
			if err != nil {
				log.Println(err)
			}
			c.msg.ReplyText(c.info.GenrateMessage(replay.SucaiDownLoadMsg))

		} else {
			c.info.ErrMsg = resp
			c.msg.ReplyText(c.info.GenrateMessage(replay.SucaiDownLoadErrMsg))
		}
	}
}

func (c Impl) UserInfo() {
	var (
		db = conn.DB
	)
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	userId, err := utils.GetUserId(c.msg)
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	var wUser model.WUser
	err = db.Model(&model.WUser{}).Where("wechat_id = ?", userId).Find(&wUser).Error
	if err != nil {
		return
	}
	c.info.UserId = strconv.Itoa(wUser.Id)
	c.info.Balance = wUser.Balance
	c.msg.ReplyText(c.info.GenrateMessage(replay.UserInfoMsg))
	return

}

func (c Impl) Recharge() {
	var (
		db = conn.DB
	)
	name, err := utils.GetUserName(c.msg)
	c.info = replay.Info{
		NickName: name,
	}
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	userId, err := utils.GetUserId(c.msg)
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	var wUser model.WUser
	err = db.Model(&model.WUser{}).Where("wechat_id = ?", userId).Find(&wUser).Error
	if err != nil {
		return
	}
	if wUser.IsAdmin == consts.IsAdminTrue {
		commands := strings.SplitN(c.realMsg, " ", 2)
		if len(commands) != 2 {
			c.msg.ReplyText(c.info.GenrateMessage(replay.TaskSubCommandErrMsg))
			return
		}
		customerId := strings.TrimSpace(commands[0])
		balanceStr := strings.ToLower(strings.TrimSpace(commands[1]))
		balance, _ := strconv.Atoi(balanceStr)
		var cUser model.WUser
		err = db.Model(&model.WUser{}).Where("id = ?", customerId).Find(&cUser).Error
		if err != nil {
			return
		}
		cUser.Balance = cUser.Balance + balance
		err = db.Save(&cUser).Error
		if err != nil {
			return
		}
		c.info.Balance = cUser.Balance
		c.info.CustomerName = cUser.NickName
		c.info.UserId = strconv.Itoa(cUser.Id)
		c.msg.ReplyText(c.info.GenrateMessage(replay.RechargeMsg))
	} else {
		c.msg.ReplyText(c.info.GenrateMessage(replay.NotAdmin))
	}

	return

}

func (c Impl) CheckUserBalance(operate int) (flag bool) {
	var (
		db = conn.DB
	)
	flag = true
	userId, err := utils.GetUserId(c.msg)
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	var wUser model.WUser
	err = db.Model(&model.WUser{}).Where("wechat_id = ?", userId).Find(&wUser).Error
	if err != nil {
		return
	}
	if (wUser.Balance - operate) < 0 {
		flag = false
		c.msg.ReplyText(c.info.GenrateMessage(replay.BalanceErr))
	}
	return
}

func (c Impl) DeductionsBalance(operate int) (err error) {
	var (
		db = conn.DB
	)
	userId, err := utils.GetUserId(c.msg)
	if err != nil {
		c.msg.ReplyText(c.info.GenrateMessage(replay.TaskNewUserErrMsg))
		return
	}
	var wUser model.WUser
	err = db.Model(&model.WUser{}).Where("wechat_id = ?", userId).Find(&wUser).Error
	if err != nil {
		return
	}
	wUser.Balance = wUser.Balance - operate
	err = db.Save(&wUser).Error
	if err != nil {
		return
	}
	return
}

func (c Impl) GetImg() (base64String string, err error) {
	// 1. 发送HTTP请求以获取网络图片
	imageURL := c.realMsg
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("无法获取图片:", err)
		return
	}
	defer response.Body.Close()
	// 2. 读取响应并将其转换为字节数组
	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("无法读取图片:", err)
		return
	}
	// 3. 使用Base64编码函数将字节数组转换为Base64编码的字符串
	base64String = base64.StdEncoding.EncodeToString(imageBytes)
	base64String = "data:image/png;base64," + base64String
	return
}

func (c Impl) GetImgs() (bsArray []string) {
	pmts := strings.Split(c.realMsg, " ")
	rexLinx := regexp.MustCompile(`https?://[\w\-.:/?&=]+`)
	for _, p := range pmts {
		if rexLinx.MatchString(p) {
			response, err := http.Get(p)
			if err != nil {
				fmt.Println("无法获取图片:", err)
				return
			}
			defer response.Body.Close()
			// 2. 读取响应并将其转换为字节数组
			imageBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("无法读取图片:", err)
				return
			}
			// 3. 使用Base64编码函数将字节数组转换为Base64编码的字符串
			base64String := base64.StdEncoding.EncodeToString(imageBytes)
			base64String = "data:image/png;base64," + base64String
			bsArray = append(bsArray, base64String)
		}
	}
	return
}

/**
欢迎使用梦幻画室为您提供的Midjourney服务
------------------------------
一、绘图功能
· 输入 /mj prompt
<prompt> 即你像mj提的绘画需求
------------------------------
二、变换功能
· 输入 /mj 1234567 U1
· 输入 /mj 1234567 V1
<1234567> 代表消息ID，<U>代表放大，<V>代表细致变化，<1>代表第几张图
------------------------------
三、附加参数
1.解释：附加参数指的是在prompt后携带的参数，可以使你的绘画更加别具一格
· 输入 /mj prompt --v 5 --ar 16:9
2.使用：需要使用--key value ，key和value之间需要空格隔开，每个附加参数之间也需要空格隔开
3.详解：上述附加参数解释 <v>版本key <5>版本号 <ar>比例key，<16:9>比例value
------------------------------
四、附加参数列表
1.(--version) 或 (--v) 《版本》 参数 1，2，3，4，5 默认4，不可与niji同用
2.(--niji)《卡通版本》 参数 空或 5 默认空，不可与版本同用
3.(--aspect) 或 (--ar) 《横纵比》 参数 n:n ，默认1:1 ,不通版本略有差异，具体详见机器人提示
4.(--chaos) 或 (--c) 《噪点》参数 0-100 默认0
5.(--quality) 或 (--q) 《清晰度》参数 .25 .5 1 2 分别代表，一般，清晰，高清，超高清，默认1
6.(--style) 《风格》参数 4a,4b,4c (v4)版本可用，参数 expressive,cute (niji5)版本可用
7.(--stylize) 或 (--s)) 《风格化》参数 1-1000 v3 625-60000
8.(--seed) 《种子》参数 0-4294967295 可自定义一个数值配合(sameseed)使用
9.(--sameseed) 《相同种子》参数 0-4294967295 可自定义一个数值配合(seed)使用
10.(--tile) 《重复模式》参数 空
*/
func (c Impl) Help() {
	msg :=
		"欢迎使用MJBOT\n" +
			"------------------------------\n" +
			"🎨 生成图片命令 \n" +
			"输入: /imagine prompt\n" +
			"<prompt> 即你向mj提的绘画需求\n" +
			"------------------------------\n" +
			"🌈 变换图片命令 ️\n" +
			"输入: /up 12345678 U1\n" +
			"输入: /up 12345678 V1\n" +
			"<12345678> 代表任务ID，<U>代表放大，<V>代表细致变化，<1>代表第几张图\n" +
			"------------------------------\n" +
			"🌈 其他图片命令 ️\n" +
			"输入: /describe 图片链接\n" +
			"输入: /blend 图片链接 图片链接 .... \n" +
			"/describe 反推关键词, /blend 融合图片[可以多张空格隔开]\n" +
			"免费上传图片网站: postimages.org 上传后复制图片地址 \n" +
			"------------------------------\n" +
			"📕 附加参数 \n" +
			"1.解释：附加参数指的是在prompt后携带的参数，可以使你的绘画更加别具一格\n" +
			"· 输入 /imagine prompt --v 5 --ar 16:9\n" +
			"2.使用：需要使用--key value ，key和value之间需要空格隔开，每个附加参数之间也需要空格隔开\n" +
			"3.详解：上述附加参数解释 <v>版本key <5>版本号 <ar>比例key，<16:9>比例value\n" +
			"------------------------------\n" +
			"📗 附加参数列表\n" +
			"1.(--version) 或 (--v) 《版本》 参数 1，2，3，4，5 默认5，不可与niji同用\n" +
			"2.(--niji)《卡通版本》 参数 空或 5 默认空，不可与版本同用\n" +
			"3.(--aspect) 或 (--ar) 《横纵比》 参数 n:n ，默认1:1 ，不同版本略有差异，具体详见机器人提示\n" +
			"4.(--chaos) 或 (--c) 《噪点》参数 0-100 默认0\n" +
			"5.(--quality) 或 (--q) 《清晰度》参数 .25 .5 1 2 分别代表，一般，清晰，高清，超高清，默认1\n" +
			"6.(--style) 《风格》参数 4a,4b,4c (v4)版本可用，参数 expressive,cute (niji5)版本可用\n" +
			"7.(--stylize) 或 (--s)) 《风格化》参数 1-1000 v3 625-60000\n" +
			"8.(--seed) 《种子》参数 0-4294967295 可自定义一个数值配合(sameseed)使用\n" +
			"9.(--sameseed) 《相同种子》参数 0-4294967295 可自定义一个数值配合(seed)使用\n" +
			"10.(--tile) 《重复模式》参数 空\n" +
			"------------------------------\n" +
			"📚 系统命令\n" +
			"/me 查询个人信息"
	c.msg.ReplyText(msg)
}

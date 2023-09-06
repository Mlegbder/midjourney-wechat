package config

type Config struct {
	Mysql struct {
		Dsn         string `json:"Dsn" yaml:"Dsn"`
		MaxIdleConn int    `json:"MaxIdleConn,default=50" yaml:"MaxIdleConn"`
		MaxOpenConn int    `json:"MaxOpenConn,default=100" yaml:"MaxOpenConn"`
	} `yaml:"Mysql"`
}

package config

import (
	mysql "application/pkg/config/mysql"
	web "application/pkg/config/web"
	"os"

	"github.com/namsral/flag"
)

type Config struct {
	IsDebug   bool        `json:"isDeubg"`
	Mysql     *mysql.Conf `json:"orm"`
	WebEngine *web.Conf   `json:"webEngine"`
}

const (
	envPrefix = "FABRIC"
)

func InitConfig() (conf *Config, err error) {
	conf = new(Config)
	fg := flag.NewFlagSetWithEnvPrefix(os.Args[0], envPrefix, flag.ContinueOnError)
	fg.String(flag.DefaultConfigFlagname, "", "配置文件路径(abs)") // 有文件解析文件
	fg.BoolVar(&conf.IsDebug, "isDebug", false, "isDebug")
	// mysql配置
	conf.Mysql, err = new(mysql.Conf).Parse(fg)
	if err != nil {
		return nil, err
	}
	// 服务配置
	conf.WebEngine, err = new(web.Conf).Parse(fg)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

package mysql

import (
	"os"

	"github.com/namsral/flag"
)

// Conf mysql的配置
type Conf struct {
	Host     string `json:"host"` // ip:port
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	DB       string `json:"db"`
}

func (conf *Conf) Parse(fg *flag.FlagSet) (*Conf, error) {
	fg.StringVar(&conf.Host, "mysql_host", "127.0.0.1:3306", "mysql的主机地址")
	fg.StringVar(&conf.Username, "mysql_user", "root", "mysql的用户名")
	fg.StringVar(&conf.Passwd, "mysql_passwd", "123456", "mysql的密码")
	fg.StringVar(&conf.DB, "mysql_DB", "fabric", "mysql的DB")
	// 执行参数解析
	if err := fg.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	return conf, nil
}

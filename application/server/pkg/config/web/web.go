package web

import (
	"os"

	"github.com/namsral/flag"
)

const (
	webHost = "0.0.0.0:8888" // 默认值: web端口
)

// Conf mysql的配置
type Conf struct {
	Host string `json:"host"` // web服务器
}

func (conf *Conf) Parse(fg *flag.FlagSet) (*Conf, error) {
	fg.StringVar(&conf.Host, "webhost", webHost, "web服务host")

	// 执行参数解析
	if err := fg.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return conf, nil
}

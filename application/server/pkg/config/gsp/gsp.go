package gsp

import (
	"os"

	"github.com/namsral/flag"
)

type Conf struct {
	AccessKey string
	SecretKey string
	Regions   string
	Addr      string
}

func (conf *Conf) Parse(fg *flag.FlagSet) (*Conf, error) {
	fg.StringVar(&conf.AccessKey, "gsp_access_key", "m0sJigHaG22k2fyIarkf", "access_key")
	fg.StringVar(&conf.SecretKey, "gsp_secret_key", "6NUCG72NDwmvoT5sUEyNeb09jJniJTJolH3s8D81", "secret_key")
	fg.StringVar(&conf.Regions, "gsp_regions", "fabric", "regions")
	fg.StringVar(&conf.Addr, "gsp_addr", "http://127.0.0.1:9000", "addr")
	// 执行参数解析
	if err := fg.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	return conf, nil
}

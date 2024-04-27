package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/config"
	"application/pkg/cron"
	orm "application/pkg/mysql"
	"application/routers"

	"gorm.io/gorm"
)

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8888)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(),
	}
	// 初始化 config 配置
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	//
	orm, err := orm.NewOrm(&orm.MysqlConfig{
		Host:     conf.Mysql.Host,
		Username: conf.Mysql.Username,
		Passwd:   conf.Mysql.Passwd,
		DB:       conf.Mysql.DB,
		IsDebug:  conf.IsDebug,
	}, new(gorm.Config))
	if err != nil {
		panic(err)
	}

	log.Printf("[info] start http server listening %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}

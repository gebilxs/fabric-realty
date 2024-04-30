package main

import (
	"context"
	"syscall"

	"application/config"
	"application/model"
	"application/pkg/gsp"
	orm "application/pkg/mysql"

	"application/services"
	"os/signal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()

	// 初始化 config 配置
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	// 增加web服务
	app := gin.Default()
	// 实例化对象存储
	gspConf := conf.GSP
	gsp, err := gsp.NewGSP(gspConf.Addr, gspConf.AccessKey, gspConf.SecretKey, gspConf.Regions)
	if err != nil {
		panic(err)
	}
	// 增加orm连接
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
	// 注册登陆连接对象
	loginorm, err := model.NewLoginTaskManager(orm)
	if err != nil {
		panic(err)
	}

	// 路线信息连接对象
	alignmentorm, err := model.NewAlignmentManager(orm)
	if err != nil {
		println("alignmentorm error")
	}

	// 音乐信息连接对象
	musicorm, err := model.NewMusicManager(orm)
	if err != nil {
		println("musicorm error")
	}
	// 车辆信息连接对象
	carorm, err := model.NewCarManager(orm)
	if err != nil {
		println("carorm error")
	}
	// 车辆行为信息连接对象
	actionorm, err := model.NewActionManager(orm)
	if err != nil {
		println("actionorm error")
	}
	// 事故信息连接对象
	accidentorm, err := model.NewAccidentManager(orm)
	if err != nil {
		println("accidentorm error")
	}

	srv := services.NewFabricSrv(ctx, conf, app, gsp, loginorm,
		alignmentorm, musicorm, carorm, actionorm, accidentorm)

	if err != srv.Start(ctx) {
		panic(err)
	}
}

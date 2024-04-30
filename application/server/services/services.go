package services

import (
	"application/blockchain"
	"application/config"
	"application/model"
	"application/pkg/cron"
	"context"
	"log"

	"application/pkg/gsp"
	"time"

	"github.com/gin-gonic/gin"
)

type FabricSrv struct {
	ctx          context.Context
	conf         *config.Config
	app          *gin.Engine
	gsp          *gsp.GSP
	loginorm     *model.LoginAndRegisterManager
	alignmentorm *model.AlignmentManager
	musicorm     *model.MusicManager
	carorm       *model.CarManager
	actionorm    *model.ActionManager
	accidentorm  *model.AccidentManager
}

type OptionFunc func(*FabricSrv)

func NewFabricSrv(ctx context.Context, conf *config.Config, app *gin.Engine, gsp *gsp.GSP,
	loginorm *model.LoginAndRegisterManager,
	alignmentorm *model.AlignmentManager, musicorm *model.MusicManager, carorm *model.CarManager,
	actionorm *model.ActionManager, accidentorm *model.AccidentManager,
	ops ...OptionFunc) *FabricSrv {
	fabsrv := &FabricSrv{
		ctx:          ctx,
		conf:         conf,
		app:          app,
		gsp:          gsp,
		loginorm:     loginorm,
		alignmentorm: alignmentorm,
		musicorm:     musicorm,
		carorm:       carorm,
		actionorm:    actionorm,
		accidentorm:  accidentorm,
	}
	for _, op := range ops {
		op(fabsrv)
	}
	return fabsrv
}

func (fab *FabricSrv) Start(ctx context.Context) error {
	signal := make(chan struct{})
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	blockchain.Init()
	go cron.Init()
	fab.InitRouter()
	// 启动应用程序的服务
	go func(subctx context.Context) {
		if err := fab.app.Run(fab.conf.WebEngine.Host); err != nil {
			println("web服务启动错误")
		}
		signal <- struct{}{}
	}(ctx)
	<-ctx.Done()
	return nil
}

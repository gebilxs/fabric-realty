package services

import (
	v1 "application/api/v1"
	"application/pkg/middleware"
)

// InitRouter 初始化路由信息
func (fab *FabricSrv) InitRouter() {
	fab.app.Use(middleware.Core())
	apiV1 := fab.app.Group("/api/v1")
	{
		// 登陆及人员管理
		apiV1.GET("/login_list", fab.GetLoginList)
		apiV1.POST("/login_add", fab.AddLogin)
		apiV1.POST("/login", fab.Login)
		apiV1.PUT("/login", fab.UpdateLogin)
		apiV1.DELETE("/login", fab.DeleteLogin)
		// 道路信息
		apiV1.GET("/alignment_list", fab.GetAlignmentList)
		apiV1.POST("/alignment", fab.AddAlignment)
		apiV1.PUT("/alignment", fab.UpdateAlignment)
		apiV1.DELETE("/alignment", fab.DeleteAlignment)
		// 音乐信息
		apiV1.GET("/music_list", fab.GetMusicList)
		apiV1.POST("/music", fab.AddMusic)
		apiV1.PUT("/music", fab.UpdateMusic)
		apiV1.DELETE("/music", fab.DeleteMusic)
		apiV1.GET("/music_download", fab.DownloadMusic)
		// 车辆信息
		apiV1.GET("/car_list", fab.GetCarList)
		apiV1.POST("/car", fab.AddCar)
		apiV1.PUT("/car", fab.UpdateCar)
		apiV1.DELETE("/car", fab.DeleteCar)
		// 行为信息
		apiV1.GET("/action_list", fab.GetActionList)
		apiV1.POST("/action", fab.AddAction)
		apiV1.PUT("/action", fab.UpdateAction)
		apiV1.DELETE("/action", fab.DeleteAction)
		// 事故信息
		apiV1.GET("/accident_list", fab.GetAccidentList)
		apiV1.POST("/accident", fab.AddAccident)
		apiV1.PUT("/accident", fab.UpdateAccident)
		apiV1.DELETE("/accident", fab.DeleteAccident)

		// ------ 区块连相关操作 ---------D
		apiV1.POST("/queryAccountList", v1.QueryAccountList)
		apiV1.POST("/createRealEstate", v1.CreateRealEstate)
		apiV1.POST("/queryRealEstateList", v1.QueryRealEstateList)
		apiV1.POST("/createSelling", v1.CreateSelling)
		apiV1.POST("/createSellingByBuy", v1.CreateSellingByBuy)
		apiV1.POST("/querySellingList", v1.QuerySellingList)
		apiV1.POST("/querySellingListByBuyer", v1.QuerySellingListByBuyer)
		apiV1.POST("/updateSelling", v1.UpdateSelling)
		apiV1.POST("/createDonating", v1.CreateDonating)
		apiV1.POST("/queryDonatingList", v1.QueryDonatingList)
		apiV1.POST("/queryDonatingListByGrantee", v1.QueryDonatingListByGrantee)
		apiV1.POST("/updateDonating", v1.UpdateDonating)
	}
}

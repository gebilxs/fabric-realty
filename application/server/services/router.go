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
		// apiV1.POST("/login_add", fab.AddLogin)
		// ------
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

package services

import (
	"application/model"
	"application/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type ReqCarList struct {
	PageIndex   int    `json:"pageIndex"`   // 当前页码,可选
	PageSize    int    `json:"pageSize"`    // 每页展示的条数. 可选
	CarName     string `json:"carName"`     // 车名
	CarType     string `json:"carType"`     // 车类型
	CarColor    string `json:"carColor"`    // 车颜色
	CarEngine   string `json:"carEngine"`   // 车引擎
	CarLocation string `json:"carLocation"` // 车位置
	CarHealth   string `json:"carHealth"`   // 车健康
}

type RespCarList struct {
	ID          uint   `json:"id"`
	CreateTime  int64  `json:"create_time"`
	CarName     string `json:"carName"`     // 车名
	CarType     string `json:"carType"`     // 车类型
	CarColor    string `json:"carColor"`    // 车颜色
	CarEngine   string `json:"carEngine"`   // 车引擎
	CarLocation string `json:"carLocation"` // 车位置
	CarHealth   string `json:"carHealth"`   // 车健康
}

// 查询车辆列表
func (fab *FabricSrv) GetCarList(ctx *gin.Context) {
	req := &ReqCarList{
		PageIndex:   cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:    cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		CarName:     ctx.DefaultQuery("carName", ""),
		CarType:     ctx.DefaultQuery("carType", ""),
		CarColor:    ctx.DefaultQuery("carColor", ""),
		CarEngine:   ctx.DefaultQuery("carEngine", ""),
		CarLocation: ctx.DefaultQuery("carLocation", ""),
		CarHealth:   ctx.DefaultQuery("carHealth", ""),
	}
	query := makeCarQuery(req)
	// 查询
	total, err := fab.carorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.carorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespCarList, 0)
	for _, item := range list {
		items = append(items, &RespCarList{
			ID:          item.ID,
			CreateTime:  item.CreatedAt.UnixMilli(),
			CarName:     item.CarName,
			CarType:     item.CarType,
			CarColor:    item.CarColor,
			CarEngine:   item.CarEngine,
			CarLocation: item.CarLocation,
			CarHealth:   item.CarHealth,
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeCarQuery(req *ReqCarList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.CarName != "" {
		query.Args = append(query.Args, req.CarName)
		if query.Query != "" {
			query.Query += " AND car_name = ? "
		} else {
			query.Query = "car_name = ? "
		}
	}
	if req.CarType != "" {
		query.Args = append(query.Args, req.CarType)
		if query.Query != "" {
			query.Query += " AND car_type = ? "
		} else {
			query.Query = "car_type = ? "
		}
	}
	if req.CarColor != "" {
		query.Args = append(query.Args, req.CarColor)
		if query.Query != "" {
			query.Query += " AND car_color = ? "
		} else {
			query.Query = "car_color = ? "
		}
	}
	if req.CarEngine != "" {
		query.Args = append(query.Args, req.CarEngine)
		if query.Query != "" {
			query.Query += " AND car_engine = ? "
		} else {
			query.Query = "car_engine = ? "
		}
	}
	if req.CarLocation != "" {
		query.Args = append(query.Args, req.CarLocation)
		if query.Query != "" {
			query.Query += " AND car_location = ? "
		} else {
			query.Query = "car_location = ? "
		}
	}
	if req.CarHealth != "" {
		query.Args = append(query.Args, req.CarHealth)
		if query.Query != "" {
			query.Query += " AND car_health = ? "
		} else {
			query.Query = "car_health = ? "
		}
	}
	return query
}

// 增加车辆信息
type ReqAddCar struct {
	CarName     string `json:"carName"`     // 车名
	CarType     string `json:"carType"`     // 车类型
	CarColor    string `json:"carColor"`    // 车颜色
	CarEngine   string `json:"carEngine"`   // 车引擎
	CarLocation string `json:"carLocation"` // 车位置
	CarHealth   string `json:"carHealth"`   // 车健康
}

func (fab *FabricSrv) AddCar(ctx *gin.Context) {
	var req ReqAddCar
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	item := &model.Car{
		CarName:     req.CarName,
		CarType:     req.CarType,
		CarColor:    req.CarColor,
		CarEngine:   req.CarEngine,
		CarLocation: req.CarLocation,
		CarHealth:   req.CarHealth,
	}
	err := fab.carorm.Insert(ctx, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("新增失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("新增成功", nil).JSON())
}

// 更新车辆信息
type ReqUpdateCar struct {
	ID          uint   `json:"id" binding:"required"`
	CarName     string `json:"carName"`     // 车名
	CarType     string `json:"carType"`     // 车类型
	CarColor    string `json:"carColor"`    // 车颜色
	CarEngine   string `json:"carEngine"`   // 车引擎
	CarLocation string `json:"carLocation"` // 车位置
	CarHealth   string `json:"carHealth"`   // 车健康
}

func (fab *FabricSrv) UpdateCar(ctx *gin.Context) {
	var req ReqUpdateCar
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	item := &model.Car{
		CarName:     req.CarName,
		CarType:     req.CarType,
		CarColor:    req.CarColor,
		CarEngine:   req.CarEngine,
		CarLocation: req.CarLocation,
		CarHealth:   req.CarHealth,
	}
	err := fab.carorm.Update(ctx, req.ID, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新成功", "").JSON())
}

// 删除车辆信息
type ReqDeleteCar struct {
	ID uint `json:"id" validate:"required"` // ID
}

func (fab *FabricSrv) DeleteCar(ctx *gin.Context) {
	var req ReqDeleteCar
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	err = fab.carorm.Delete(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除成功", "").JSON())
}

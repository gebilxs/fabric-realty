package services

import (
	"application/model"
	"application/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type ReqActionList struct {
	PageIndex    int    `json:"pageIndex"` // 当前页码,可选
	PageSize     int    `json:"pageSize"`  // 每页展示的条数. 可选
	ActionName   string `json:"actionName"`
	ActionType   string `json:"actionType"`
	Driver       string `json:"driver"`
	Acceleration string `json:"acceleration"`
	Steering     string `json:"steering"`
	Result       string `json:"result"`
}

type RespActionList struct {
	ID           uint   `json:"id"`
	ActionName   string `json:"actionName"`
	ActionType   string `json:"actionType"`
	Driver       string `json:"driver"`
	Acceleration string `json:"acceleration"`
	Steering     string `json:"steering"`
	Result       string `json:"result"`
}

func (fab *FabricSrv) GetActionList(ctx *gin.Context) {
	req := &ReqActionList{
		PageIndex:    cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:     cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		ActionName:   ctx.DefaultQuery("actionName", ""),
		ActionType:   ctx.DefaultQuery("actionType", ""),
		Driver:       ctx.DefaultQuery("driver", ""),
		Acceleration: ctx.DefaultQuery("acceleration", ""),
		Steering:     ctx.DefaultQuery("steering", ""),
		Result:       ctx.DefaultQuery("result", ""),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	query := makeActionQuery(req)
	total, err := fab.actionorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.actionorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespActionList, 0)
	for _, item := range list {
		items = append(items, &RespActionList{
			ID:           item.ID,
			ActionName:   item.ActionName,
			ActionType:   item.ActionType,
			Driver:       item.Driver,
			Acceleration: item.Acceleration,
			Steering:     item.Steering,
			Result:       item.Result,
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeActionQuery(req *ReqActionList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.ActionName != "" {
		query.Args = append(query.Args, req.ActionName)
		if query.Query != "" {
			query.Query += " AND action_name = ? "
		} else {
			query.Query = "action_name = ? "
		}
	}
	if req.ActionType != "" {
		query.Args = append(query.Args, req.ActionType)
		if query.Query != "" {
			query.Query += " AND action_type = ? "
		} else {
			query.Query = "action_type = ? "
		}
	}
	if req.Driver != "" {
		query.Args = append(query.Args, req.Driver)
		if query.Query != "" {
			query.Query += " AND driver = ? "
		} else {
			query.Query = "driver = ? "
		}
	}
	if req.Acceleration != "" {
		query.Args = append(query.Args, req.Acceleration)
		if query.Query != "" {
			query.Query += " AND acceleration = ? "
		} else {
			query.Query = "acceleration = ? "
		}
	}
	if req.Steering != "" {
		query.Args = append(query.Args, req.Steering)
		if query.Query != "" {
			query.Query += " AND steering = ? "
		} else {
			query.Query = "steering = ? "
		}
	}
	if req.Result != "" {
		query.Args = append(query.Args, req.Result)
		if query.Query != "" {
			query.Query += " AND result = ? "
		} else {
			query.Query = "result = ? "
		}
	}

	return query
}

// 新增行为分析
type ReqAddAction struct {
	ActionName   string `json:"actionName"`
	ActionType   string `json:"actionType"`
	Driver       string `json:"driver"`
	Acceleration string `json:"acceleration"`
	Steering     string `json:"steering"`
	Result       string `json:"result"`
}

func (fab *FabricSrv) AddAction(ctx *gin.Context) {
	req := &ReqAddAction{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	action := &model.Action{
		ActionName:   req.ActionName,
		ActionType:   req.ActionType,
		Driver:       req.Driver,
		Acceleration: req.Acceleration,
		Steering:     req.Steering,
		Result:       req.Result,
	}
	if err := fab.actionorm.Insert(ctx, action); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("新增失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("新增成功", nil).JSON())
}

// 更新行为分析
type ReqUpdateAction struct {
	ID           uint   `json:"id"`
	ActionName   string `json:"actionName"`
	ActionType   string `json:"actionType"`
	Driver       string `json:"driver"`
	Acceleration string `json:"acceleration"`
	Steering     string `json:"steering"`
	Result       string `json:"result"`
}

func (fab *FabricSrv) UpdateAction(ctx *gin.Context) {
	req := &ReqUpdateAction{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	action := &model.Action{
		ActionName:   req.ActionName,
		ActionType:   req.ActionType,
		Driver:       req.Driver,
		Acceleration: req.Acceleration,
		Steering:     req.Steering,
		Result:       req.Result,
	}
	if err := fab.actionorm.Update(ctx, req.ID, action); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新成功", "").JSON())
}

// 删除行为分析
type ReqDeleteAction struct {
	ID uint `json:"id" validate:"required"` // ID
}

func (fab *FabricSrv) DeleteAction(ctx *gin.Context) {
	req := &ReqDeleteAction{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := fab.actionorm.Delete(ctx, req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除成功", "").JSON())
}

package services

import (
	"application/model"
	"application/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type ReqAccidentList struct {
	PageIndex        int    `json:"pageIndex"` // 当前页码,可选
	PageSize         int    `json:"pageSize"`  // 每页展示的条数. 可选
	AccidentName     string `json:"accidentName"`
	AccidentType     string `json:"accidentType"`
	Driver           string `json:"driver"`
	AccidentLocation string `json:"accidentLocation"`
	AccidentNotice   string `json:"accidentNotice"`
	AccidentResult   string `json:"accidentResult"`
	AccidentReason   string `json:"accidentReason"`
	AccidentStatus   string `json:"accidentStatus"`
}

type RespAccidentList struct {
	ID               uint   `json:"id"`
	CreateTime       int64  `json:"create_time"`
	AccidentName     string `json:"accidentName"`
	AccidentType     string `json:"accidentType"`
	Driver           string `json:"driver"`
	AccidentLocation string `json:"accidentLocation"`
	AccidentNotice   string `json:"accidentNotice"`
	AccidentResult   string `json:"accidentResult"`
	AccidentReason   string `json:"accidentReason"`
	AccidentStatus   string `json:"accidentStatus"`
}

// 分页查询
func (fab *FabricSrv) GetAccidentList(ctx *gin.Context) {
	req := &ReqAccidentList{
		PageIndex:        cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:         cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		AccidentName:     ctx.Query("accidentName"),
		AccidentType:     ctx.Query("accidentType"),
		Driver:           ctx.Query("driver"),
		AccidentLocation: ctx.Query("accidentLocation"),
		AccidentNotice:   ctx.Query("accidentNotice"),
		AccidentResult:   ctx.Query("accidentResult"),
		AccidentReason:   ctx.Query("accidentReason"),
		AccidentStatus:   ctx.Query("accidentStatus"),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	query := makeAccidentQuery(req)
	total, err := fab.accidentorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.accidentorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespAccidentList, 0)
	for _, item := range list {
		items = append(items, &RespAccidentList{
			ID:               item.ID,
			CreateTime:       item.CreatedAt.UnixMilli(),
			AccidentName:     item.AccidentName,
			AccidentType:     item.AccidentType,
			Driver:           item.Driver,
			AccidentLocation: item.AccidentLocation,
			AccidentNotice:   item.AccidentNotice,
			AccidentResult:   item.AccidentResult,
			AccidentReason:   item.AccidentReason,
			AccidentStatus:   item.AccidentStatus,
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeAccidentQuery(req *ReqAccidentList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.AccidentName != "" {
		query.Args = append(query.Args, req.AccidentName)
		if query.Query != "" {
			query.Query += " AND accident_name = ? "
		} else {
			query.Query = "accident_name = ? "
		}
	}
	if req.AccidentType != "" {
		query.Args = append(query.Args, req.AccidentType)
		if query.Query != "" {
			query.Query += " AND accident_type = ? "
		} else {
			query.Query = "accident_type = ? "
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
	if req.AccidentLocation != "" {
		query.Args = append(query.Args, req.AccidentLocation)
		if query.Query != "" {
			query.Query += " AND accident_location = ? "
		} else {
			query.Query = "accident_location = ? "
		}
	}
	if req.AccidentNotice != "" {
		query.Args = append(query.Args, req.AccidentNotice)
		if query.Query != "" {
			query.Query += " AND accident_notice = ? "
		} else {
			query.Query = "accident_notice = ? "
		}
	}
	if req.AccidentResult != "" {
		query.Args = append(query.Args, req.AccidentResult)
		if query.Query != "" {
			query.Query += " AND accident_result = ? "
		} else {
			query.Query = "accident_result = ? "
		}
	}
	if req.AccidentReason != "" {
		query.Args = append(query.Args, req.AccidentReason)
		if query.Query != "" {
			query.Query += " AND accident_reason = ? "
		} else {
			query.Query = "accident_reason = ? "
		}
	}
	if req.AccidentStatus != "" {
		query.Args = append(query.Args, req.AccidentStatus)
		if query.Query != "" {
			query.Query += " AND accident_status = ? "
		} else {
			query.Query = "accident_status = ? "
		}
	}
	return query
}

type ReqAddAccident struct {
	AccidentName     string `json:"accidentName"`
	AccidentType     string `json:"accidentType"`
	Driver           string `json:"driver"`
	AccidentLocation string `json:"accidentLocation"`
	AccidentNotice   string `json:"accidentNotice"`
	AccidentResult   string `json:"accidentResult"`
	AccidentReason   string `json:"accidentReason"`
	AccidentStatus   string `json:"accidentStatus"`
}

func (fab *FabricSrv) AddAccident(ctx *gin.Context) {
	var req ReqAddAccident
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	accident := &model.Accident{
		AccidentName:     req.AccidentName,
		AccidentType:     req.AccidentType,
		Driver:           req.Driver,
		AccidentLocation: req.AccidentLocation,
		AccidentNotice:   req.AccidentNotice,
		AccidentResult:   req.AccidentResult,
		AccidentReason:   req.AccidentReason,
		AccidentStatus:   req.AccidentStatus,
	}
	if err := fab.accidentorm.Insert(ctx, accident); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("添加失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("添加成功", nil).JSON())
}

type ReqUpdateAccident struct {
	ID               uint   `json:"id"`
	AccidentName     string `json:"accidentName"`
	AccidentType     string `json:"accidentType"`
	Driver           string `json:"driver"`
	AccidentLocation string `json:"accidentLocation"`
	AccidentNotice   string `json:"accidentNotice"`
	AccidentResult   string `json:"accidentResult"`
	AccidentReason   string `json:"accidentReason"`
	AccidentStatus   string `json:"accidentStatus"`
}

func (fab *FabricSrv) UpdateAccident(ctx *gin.Context) {
	var req ReqUpdateAccident
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	accident := &model.Accident{
		AccidentName:     req.AccidentName,
		AccidentType:     req.AccidentType,
		Driver:           req.Driver,
		AccidentLocation: req.AccidentLocation,
		AccidentNotice:   req.AccidentNotice,
		AccidentResult:   req.AccidentResult,
		AccidentReason:   req.AccidentReason,
		AccidentStatus:   req.AccidentStatus,
	}
	if err := fab.accidentorm.Update(ctx, req.ID, accident); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新成功", nil).JSON())
}

type ReqDeleteAccident struct {
	ID uint `json:"id"`
}

func (fab *FabricSrv) DeleteAccident(ctx *gin.Context) {
	var req ReqDeleteAccident
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := fab.accidentorm.Delete(ctx, req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除成功", nil).JSON())
}

package services

import (
	"application/model"
	"application/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type ReqAlignmentList struct {
	PageIndex     int    `json:"pageIndex"` // 当前页码,可选
	PageSize      int    `json:"pageSize"`  // 每页展示的条数. 可选
	RoadName      string `json:"road_name"`
	RoadType      string `json:"road_type"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	RoadDetail    string `json:"road_detail"`
	RoadStatus    string `json:"road_status"`
	BeLongitudes  string `json:"be_longitudes"`
	BeLatitudes   string `json:"be_latitudes"`
	EndLongitudes string `json:"end_longitudes"`
	EndLatitudes  string `json:"end_latitudes"`
}

type RespAlignmentList struct {
	ID            uint   `json:"id"`
	RoadName      string `json:"road_name"`
	RoadType      string `json:"road_type"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	RoadDetail    string `json:"road_detail"`
	RoadStatus    string `json:"road_status"`
	BeLongitudes  string `json:"be_longitudes"`
	BeLatitudes   string `json:"be_latitudes"`
	EndLongitudes string `json:"end_longitudes"`
	EndLatitudes  string `json:"end_latitudes"`
	CreateTime    int64  `json:"create_time"`
}

func (fab *FabricSrv) GetAlignmentList(ctx *gin.Context) {
	req := &ReqAlignmentList{
		PageIndex:     cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:      cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		RoadName:      ctx.DefaultQuery("road_name", ""),
		RoadType:      ctx.DefaultQuery("road_type", ""),
		Length:        cast.ToInt(ctx.DefaultQuery("length", "0")),
		Width:         cast.ToInt(ctx.DefaultQuery("width", "0")),
		RoadDetail:    ctx.DefaultQuery("road_detail", ""),
		RoadStatus:    ctx.DefaultQuery("road_status", ""),
		BeLongitudes:  ctx.DefaultQuery("be_longitudes", ""),
		BeLatitudes:   ctx.DefaultQuery("be_latitudes", ""),
		EndLongitudes: ctx.DefaultQuery("end_longitudes", ""),
		EndLatitudes:  ctx.DefaultQuery("end_latitudes", ""),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	query := makeAlignmentQuery(req)

	total, err := fab.alignmentorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.alignmentorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespAlignmentList, 0)
	for _, item := range list {
		items = append(items, &RespAlignmentList{
			ID:            item.ID,
			RoadName:      item.RoadName,
			RoadType:      item.RoadType,
			Length:        item.Length,
			Width:         item.Width,
			RoadDetail:    item.RoadDetail,
			RoadStatus:    item.RoadStatus,
			BeLongitudes:  item.BeLongitudes,
			BeLatitudes:   item.BeLatitudes,
			EndLongitudes: item.EndLongitudes,
			EndLatitudes:  item.EndLatitudes,
			CreateTime:    item.CreatedAt.UnixMilli(),
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeAlignmentQuery(req *ReqAlignmentList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.RoadName != "" {
		query.Args = append(query.Args, req.RoadName)
		if query.Query != "" {
			query.Query += " AND road_name = ? "
		} else {
			query.Query = "road_name = ? "
		}
	}
	if req.RoadType != "" {
		query.Args = append(query.Args, req.RoadType)
		if query.Query != "" {
			query.Query += " AND road_type = ? "
		} else {
			query.Query = "road_type = ? "
		}
	}
	if req.Length != 0 {
		query.Args = append(query.Args, req.Length)
		if query.Query != "" {
			query.Query += " AND length = ? "
		} else {
			query.Query = "length = ? "
		}
	}
	if req.Width != 0 {
		query.Args = append(query.Args, req.Width)
		if query.Query != "" {
			query.Query += " AND width = ? "
		} else {
			query.Query = "width = ? "
		}
	}
	if req.RoadDetail != "" {
		query.Args = append(query.Args, req.RoadDetail)
		if query.Query != "" {
			query.Query += " AND road_detail = ? "
		} else {
			query.Query = "road_detail = ? "
		}
	}
	if req.RoadStatus != "" {
		query.Args = append(query.Args, req.RoadStatus)
		if query.Query != "" {
			query.Query += " AND road_status = ? "
		} else {
			query.Query = "road_status = ? "
		}
	}
	if req.BeLongitudes != "" {
		query.Args = append(query.Args, req.BeLongitudes)
		if query.Query != "" {
			query.Query += " AND be_longitudes = ? "
		} else {
			query.Query = "be_longitudes = ? "
		}
	}
	if req.BeLatitudes != "" {
		query.Args = append(query.Args, req.BeLatitudes)
		if query.Query != "" {
			query.Query += " AND be_latitudes = ? "
		} else {
			query.Query = "be_latitudes = ? "
		}
	}
	if req.EndLongitudes != "" {
		query.Args = append(query.Args, req.EndLongitudes)
		if query.Query != "" {
			query.Query += " AND end_longitudes = ? "
		} else {
			query.Query = "end_longitudes = ? "
		}
	}
	if req.EndLatitudes != "" {
		query.Args = append(query.Args, req.EndLatitudes)
		if query.Query != "" {
			query.Query += " AND end_latitudes = ? "
		} else {
			query.Query = "end_latitudes = ? "
		}
	}
	return query
}

// AddAlignment 增加道路情况

type ReqAddAlignment struct {
	RoadName      string `json:"road_name"`
	RoadType      string `json:"road_type"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	RoadDetail    string `json:"road_detail"`
	RoadStatus    string `json:"road_status"`
	BeLongitudes  string `json:"be_longitudes"`
	BeLatitudes   string `json:"be_latitudes"`
	EndLongitudes string `json:"end_longitudes"`
	EndLatitudes  string `json:"end_latitudes"`
}

func (fab *FabricSrv) AddAlignment(ctx *gin.Context) {
	var req ReqAddAlignment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	item := &model.Alignment{
		RoadName:      req.RoadName,
		RoadType:      req.RoadType,
		Length:        req.Length,
		Width:         req.Width,
		RoadDetail:    req.RoadDetail,
		RoadStatus:    req.RoadStatus,
		BeLongitudes:  req.BeLongitudes,
		BeLatitudes:   req.BeLatitudes,
		EndLongitudes: req.EndLongitudes,
		EndLatitudes:  req.EndLatitudes,
	}
	err := fab.alignmentorm.Insert(ctx, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("增加道路情况失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("增加道路情况成功", nil).JSON())
}

// 修改道路情况
type ReqUpdateAlignment struct {
	ID            uint   `json:"id"`
	RoadName      string `json:"road_name"`
	RoadType      string `json:"road_type"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	RoadDetail    string `json:"road_detail"`
	RoadStatus    string `json:"road_status"`
	BeLongitudes  string `json:"be_longitudes"`
	BeLatitudes   string `json:"be_latitudes"`
	EndLongitudes string `json:"end_longitudes"`
	EndLatitudes  string `json:"end_latitudes"`
}

func (fab *FabricSrv) UpdateAlignment(ctx *gin.Context) {
	var req ReqUpdateAlignment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	item := &model.Alignment{
		RoadName:      req.RoadName,
		RoadType:      req.RoadType,
		Length:        req.Length,
		Width:         req.Width,
		RoadDetail:    req.RoadDetail,
		RoadStatus:    req.RoadStatus,
		BeLongitudes:  req.BeLongitudes,
		BeLatitudes:   req.BeLatitudes,
		EndLongitudes: req.EndLongitudes,
		EndLatitudes:  req.EndLatitudes,
	}
	err := fab.alignmentorm.Update(ctx, req.ID, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新道路情况失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新道路情况成功", nil).JSON())
}

// 删除道路情况
type ReqDeleteAlignment struct {
	ID uint `json:"id"`
}

func (fab *FabricSrv) DeleteAlignment(ctx *gin.Context) {
	var req ReqDeleteAlignment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	err := fab.alignmentorm.Delete(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除道路情况失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除道路情况成功", nil).JSON())
}

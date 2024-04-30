package services

import (
	"application/model"
	"application/pkg/app"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type ReqGetMusicList struct {
	PageIndex   int    `json:"pageIndex"` // 当前页码,可选
	PageSize    int    `json:"pageSize"`  // 每页展示的条数. 可选
	MusicName   string `json:"music_name"`
	MusicKey    string `json:"music_key"`
	MusicAuthor string `json:"music_author"`
	MusicType   string `json:"music_type"`
	MusicAlbum  string `json:"music_album"`
}

type RespMusicList struct {
	MusicName   string `json:"music_name"`
	MusicKey    string `json:"music_key"`
	MusicAuthor string `json:"music_author"`
	MusicType   string `json:"music_type"`
	MusicAlbum  string `json:"music_album"`
	CreateTime  int64  `json:"create_time"`
}

func (fab *FabricSrv) GetMusicList(ctx *gin.Context) {
	req := &ReqGetMusicList{
		PageIndex:  cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:   cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		MusicName:  ctx.DefaultQuery("music_name", ""),
		MusicKey:   ctx.DefaultQuery("music_key", ""),
		MusicType:  ctx.DefaultQuery("music_type", ""),
		MusicAlbum: ctx.DefaultQuery("music_album", ""),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	query := makeMusicQuery(req)
	total, err := fab.musicorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.musicorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespMusicList, 0)
	for _, item := range list {
		items = append(items, &RespMusicList{
			MusicName:   item.MusicName,
			MusicKey:    item.MusicKey,
			MusicAuthor: item.MusicAuthor,
			MusicType:   item.MusicType,
			MusicAlbum:  item.MusicAlbum,
			CreateTime:  item.CreatedAt.UnixMilli(),
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeMusicQuery(req *ReqGetMusicList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.MusicName != "" {
		query.Args = append(query.Args, req.MusicName)
		if query.Query != "" {
			query.Query += " AND music_name = ? "
		} else {
			query.Query = "music_name = ? "
		}
	}
	if req.MusicKey != "" {
		query.Args = append(query.Args, req.MusicKey)
		if query.Query != "" {
			query.Query += " AND music_key = ? "
		} else {
			query.Query = "music_key = ? "
		}
	}
	if req.MusicType != "" {
		query.Args = append(query.Args, req.MusicType)
		if query.Query != "" {
			query.Query += " AND music_type = ? "
		} else {
			query.Query = "music_type = ? "
		}
	}
	if req.MusicAlbum != "" {
		query.Args = append(query.Args, req.MusicAlbum)
		if query.Query != "" {
			query.Query += " AND music_album = ? "
		} else {
			query.Query = "music_album = ? "
		}
	}
	if req.MusicAuthor != "" {
		query.Args = append(query.Args, req.MusicAuthor)
		if query.Query != "" {
			query.Query += " AND music_author = ? "
		} else {
			query.Query = "music_author = ? "
		}
	}
	return query
}

// 新增音乐
type ReqAddMusic struct {
	MusicName    string                `form:"music_name" json:"music_name" binding:"required"`
	MusicKey     string                `form:"music_key" json:"music_key"`
	MusicAuthor  string                `form:"music_author" json:"music_author"`
	MusicType    string                `form:"music_type" json:"music_type"`
	MusicAlbum   string                `form:"music_album" json:"music_album"`
	MusicContent *multipart.FileHeader `form:"music_content" json:"music_content" binding:"required"`
}

func (fab *FabricSrv) AddMusic(ctx *gin.Context) {
	req := &ReqAddMusic{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	// 打开文件
	file, err := req.MusicContent.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	_, err = fab.gsp.PutS3Object(ctx, "fabric", req.MusicKey, "application/octet-stream", data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("上传文件到 S3 失败", err.Error()).JSON())
		return
	}
	item := &model.Music{
		MusicName:   req.MusicName,
		MusicKey:    req.MusicKey,
		MusicAuthor: req.MusicAuthor,
		MusicType:   req.MusicType,
		MusicAlbum:  req.MusicAlbum,
	}
	if err = fab.musicorm.Insert(ctx, item); err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("新增音乐失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("新增音乐成功", nil).JSON())

}

// 更新音乐信息
type ReqUpdateMusic struct {
	ID          uint   `json:"id" binding:"required"`
	MusicName   string `form:"music_name" json:"music_name"`
	MusicKey    string `form:"music_key" json:"music_key"`
	MusicAuthor string `form:"music_author" json:"music_author"`
	MusicType   string `form:"music_type" json:"music_type"`
	MusicAlbum  string `form:"music_album" json:"music_album"`
}

func (fab *FabricSrv) UpdateMusic(ctx *gin.Context) {
	var req ReqUpdateMusic
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	item := &model.Music{
		MusicName:   req.MusicName,
		MusicKey:    req.MusicKey,
		MusicAuthor: req.MusicAuthor,
		MusicType:   req.MusicType,
		MusicAlbum:  req.MusicAlbum,
	}
	err := fab.musicorm.Update(ctx, req.ID, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新成功", "").JSON())
}

// 删除音乐信息
type ReqDeleteMusic struct {
	ID uint `json:"id" validate:"required"` // ID
}

func (fab *FabricSrv) DeleteMusic(ctx *gin.Context) {
	var req ReqDeleteMusic
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if req.ID == 0 {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", "id不能为空").JSON())
		return
	}
	err = fab.musicorm.Delete(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除成功", "").JSON())
}

// 下载音乐信息
type ReqDownloadMusic struct {
	MusicKey string `json:"music_key" validate:"required"` // ID
}

func (fab *FabricSrv) DownloadMusic(ctx *gin.Context) {
	req := &ReqDownloadMusic{
		MusicKey: ctx.DefaultQuery("music_key", ""),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	data, err := fab.gsp.GetS3Object(ctx, "fabric", req.MusicKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("下载失败", err.Error()).JSON())
		return
	}
	ctx.Data(http.StatusOK, "application/octet-stream", data)
}

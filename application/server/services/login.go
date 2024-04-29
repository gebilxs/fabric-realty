package services

import (
	"application/pkg/app"
	"net/http"

	"application/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    "",
		"message": "string",
		"ok":      true,
	})
}

type ReqLoginList struct {
	PageIndex int    `json:"pageIndex"  validate:"required"` // 当前页码,可选
	PageSize  int    `json:"pageSize"  validate:"required"`  // 每页展示的条数. 可选
	NickName  string `json:"nickName"`                       // 昵称
	Phone     string `json:"phone"`                          // 手机号
	Email     string `json:"email"`                          // 邮箱
	Sex       string `json:"sex"`                            // 性别
	Address   string `json:"address"`                        // 地址
	Age       int    `json:"age"`                            // 年龄
}

type RespLoginList struct {
	ID         uint   `json:"id"`         // ID
	CreateTime int64  `json:"createTime"` // 创建时间
	NickName   string `json:"nickName"`   // 昵称
	Phone      string `json:"phone"`      // 手机号
	Email      string `json:"email"`      // 邮箱
	Sex        string `json:"sex"`        // 性别
	Address    string `json:"address"`    // 地址
	Age        int    `json:"age"`        // 年龄
}

// 分页查询
func (fab *FabricSrv) GetLoginList(ctx *gin.Context) {
	req := &ReqLoginList{
		PageIndex: cast.ToInt(ctx.DefaultQuery("pageIndex", "1")),
		PageSize:  cast.ToInt(ctx.DefaultQuery("pageSize", "10")),
		NickName:  ctx.Query("nickName"),
		Phone:     ctx.Query("phone"),
		Email:     ctx.Query("email"),
		Sex:       ctx.Query("sex"),
		Address:   ctx.Query("address"),
		Age:       cast.ToInt(ctx.Query("age")),
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	query := makeLoginQuery(req)
	// 查询
	total, err := fab.loginorm.QueryCount(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	list, err := fab.loginorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	items := make([]*RespLoginList, 0)
	for _, item := range list {
		items = append(items, &RespLoginList{
			ID:         item.ID,
			CreateTime: item.CreatedAt.UnixMilli(),
			NickName:   item.NickName,
			Phone:      item.Phone,
			Email:      item.Email,
			Sex:        item.Sex,
			Address:    item.Address,
			Age:        item.Age,
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeLoginQuery(req *ReqLoginList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.NickName != "" {
		query.Args = append(query.Args, req.NickName)
		if query.Query != "" {
			query.Query += " AND nick_name = ? "
		} else {
			query.Query = "nick_name = ? "
		}
	}
	if req.Phone != "" {
		query.Args = append(query.Args, req.Phone)
		if query.Query != "" {
			query.Query += " AND phone = ? "
		} else {
			query.Query = "phone = ? "
		}
	}
	if req.Email != "" {
		query.Args = append(query.Args, req.Email)
		if query.Query != "" {
			query.Query += " AND email = ? "
		} else {
			query.Query = "email = ? "
		}
	}
	if req.Sex != "" {
		query.Args = append(query.Args, req.Sex)
		if query.Query != "" {
			query.Query += " AND sex = ? "
		} else {
			query.Query = "sex = ? "
		}
	}
	if req.Address != "" {
		query.Args = append(query.Args, req.Address)
		if query.Query != "" {
			query.Query += " AND address = ? "
		} else {
			query.Query = "address = ? "
		}
	}
	if req.Age > 0 {
		query.Args = append(query.Args, req.Age)
		if query.Query != "" {
			query.Query += " AND age = ? "
		} else {
			query.Query = "age = ? "
		}
	}
	return query
}

// 增加人员

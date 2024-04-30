package services

import (
	"application/pkg/app"
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"net/http"

	"application/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

// 登陆接口 - 如果登陆成功则返回0 否则为其他
func (fab *FabricSrv) Login(ctx *gin.Context) {
	var req ReqLoginList
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	// 查询是否存在username
	query := makeLoginQuery(&req)
	list, err := fab.loginorm.QueryList(ctx, query, model.QueryWithPage(req.PageIndex, req.PageSize))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("查询失败", err.Error()).JSON())
		return
	}
	if len(list) == 0 {
		ctx.JSON(http.StatusOK, app.ResponseError("查询成功", "用户不存在").JSON())
		return
	}
	encryptedPassword, err := encryptDES([]byte(req.PassWord), []byte(DESKEY))
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return
	} else {
		ctx.JSON(http.StatusOK, app.ResponseError("查询成功", "密码错误").JSON())
	}
	if encryptedPassword == list[0].Password {
		ctx.JSON(http.StatusOK, app.ResponseOK("查询成功", "登陆成功").JSON())
		return
	}
}

type ReqLoginList struct {
	PageIndex int    `json:"pageIndex"` // 当前页码,可选
	PageSize  int    `json:"pageSize"`  // 每页展示的条数. 可选
	Username  string `json:"username"`  // 用户名
	PassWord  string `json:"password"`  // 密码
	NickName  string `json:"nickName"`  // 昵称
	Phone     string `json:"phone"`     // 手机号
	Email     string `json:"email"`     // 邮箱
	Sex       string `json:"sex"`       // 性别
	Address   string `json:"address"`   // 地址
	Age       int    `json:"age"`       // 年龄
}

type RespLoginList struct {
	ID         uint   `json:"id"`         // ID
	CreateTime int64  `json:"createTime"` // 创建时间
	Username   string `json:"username"`   // 用户名
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
			Username:   item.Username,
		})
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("查询任务成功", &app.ResponseList{
		Total: cast.ToInt(total),
		List:  items,
	}).JSON())
}

func makeLoginQuery(req *ReqLoginList) *model.WhereQuery {
	query := &model.WhereQuery{Query: "", Args: []any{}}
	if req.Username != "" {
		query.Args = append(query.Args, req.Username)
		if query.Query != "" {
			query.Query += " AND username = ? "
		} else {
			query.Query = "username = ? "
		}
	}
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

type ReqAddLogin struct {
	Username string `json:"username" validate:"required"` // 用户名
	Password string `json:"password" validate:"required"` // 密码
	NickName string `json:"nickName"`                     // 昵称
	Phone    string `json:"phone"`                        // 手机号
	Email    string `json:"email"`                        // 邮箱
	Sex      string `json:"sex`                           // 性别
	Address  string `json:"address"`                      // 地址
	Age      int    `json:"age"`                          // 年龄
}

func (fab *FabricSrv) AddLogin(ctx *gin.Context) {
	var req ReqAddLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	// 查询是否存在
	encryptedPassword, err := encryptDES([]byte(req.Password), []byte(DESKEY))
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return
	}
	item := &model.LoginAndRegister{
		Username: req.Username,
		Password: encryptedPassword,
		NickName: req.NickName,
		Phone:    req.Phone,
		Email:    req.Email,
		Sex:      req.Sex,
		Address:  req.Address,
		Age:      req.Age,
	}

	err = fab.loginorm.Insert(ctx, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("插入失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("插入成功", "").JSON())
}

// 更新账号
type ReqUpdateLogin struct {
	ID       uint   `json:"id" validate:"required"`       // ID
	Username string `json:"username" validate:"required"` // 用户名
	Password string `json:"password" validate:"required"` // 密码
	NickName string `json:"nickName"`                     // 昵称
	Phone    string `json:"phone"`                        // 手机号
	Email    string `json:"email"`                        // 邮箱
	Sex      string `json:"sex`                           // 性别
	Address  string `json:"address"`                      // 地址
	Age      int    `json:"age"`                          // 年龄
}

func (fab *FabricSrv) UpdateLogin(ctx *gin.Context) {
	var req ReqUpdateLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if err := validator.New().Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	// 查询是否存在
	encryptedPassword, err := encryptDES([]byte(req.Password), []byte(DESKEY))
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return
	}
	item := &model.LoginAndRegister{
		Username: req.Username,
		Password: encryptedPassword,
		NickName: req.NickName,
		Phone:    req.Phone,
		Email:    req.Email,
		Sex:      req.Sex,
		Address:  req.Address,
		Age:      req.Age,
	}
	err = fab.loginorm.Update(ctx, req.ID, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("更新失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("更新成功", "").JSON())
}

// 删除账号
type ReqDeleteLogin struct {
	ID uint `json:"id" validate:"required"` // ID
}

func (fab *FabricSrv) DeleteLogin(ctx *gin.Context) {
	var req ReqDeleteLogin
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", err.Error()).JSON())
		return
	}
	if req.ID == 0 {
		ctx.JSON(http.StatusBadRequest, app.ResponseError("参数校验错误", "id不能为空").JSON())
		return
	}
	err = fab.loginorm.Delete(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, app.ResponseError("删除失败", err.Error()).JSON())
		return
	}
	ctx.JSON(http.StatusOK, app.ResponseOK("删除成功", "").JSON())
}

// 密码加密部分
const DESKEY = "xck66678"

func encryptDES(plaintext, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充明文以满足DES的块大小
	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	ciphertext := make([]byte, len(plaintext))
	// 使用ECB模式加密
	ecb := cipher.NewCBCEncrypter(block, key)
	ecb.CryptBlocks(ciphertext, plaintext)

	// 将加密后的密文转换为Base64字符串
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// PKCS5Padding 填充明文以满足块大小
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

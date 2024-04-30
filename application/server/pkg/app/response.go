package app

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode int, errMsg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: 0,
		Msg:  errMsg,
		Data: data,
	})
	return
}

func (resp *Response) JSON() gin.H {
	dataB, _ := json.Marshal(resp)
	h := new(gin.H)
	_ = json.Unmarshal(dataB, h)
	return *h
}

func ResponseOK(msg string, data any) *Response {
	return &Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}

func ResponseError(msg string, data any) *Response {
	return &Response{
		Code: -1,
		Msg:  msg,
		Data: data,
	}
}

type ResponseList struct {
	Total int `json:"total"`
	List  any `json:"list"`
}

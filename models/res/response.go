package res

import (
	CODE "Smart_delivery_locker/models/res/code"
	"Smart_delivery_locker/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 封装响应
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type ListResponse[T any] struct {
	Count int64 `json:"count"`
	List  T     `json:"list"`
}

type ListMsgResponse[T any] struct {
	Count   int64  `json:"count"`
	ListMsg string `json:"list_msg"`
	List    any    `json:"list"`
}

const (
	SUCCESS = 0
	ERR     = 7
)

func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// ResultOK 成功响应
func ResultOK(data any, msg string, c *gin.Context) {
	Result(SUCCESS, data, msg, c)
}

func ResultOkWithData(data any, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

func ResultOkWithList(list any, count int64, c *gin.Context) {
	ResultOkWithData(ListResponse[any]{
		Count: count,
		List:  list,
	}, c)
}

func ResultOkWithListMsg(list any, count int64, msg string, c *gin.Context) {
	ResultOkWithData(ListMsgResponse[any]{
		Count:   count,
		ListMsg: msg,
		List:    list,
	}, c)
}

func ResultOkWithMsg(msg string, c *gin.Context) {
	Result(SUCCESS, map[string]any{}, msg, c)
}

// ResultFail 失败响应
func ResultFail(data any, msg string, c *gin.Context) {
	Result(ERR, data, msg, c)
}

func ResultFailWithMsg(msg string, c *gin.Context) {
	Result(ERR, map[string]any{}, msg, c)
}

func ResultFailWithError(err error, obj any, c *gin.Context) {
	msg := utils.GetValidMsg(err, obj)
	ResultFailWithMsg(msg, c)
}

func ResultFailWithCode(code CODE.ErrorCode, c *gin.Context) {
	if msg, ok := CODE.ErrorMap[code]; ok {
		Result(int(code), map[string]any{}, msg, c)
		return
	}

	Result(ERR, map[string]any{}, "未知错误", c)

}

package resp

import (
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ArrayResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Total   int64  `json:"total"`
	Data    any    `json:"data"`
}

func response(c *gin.Context, code int, errCode int, msg string, data any, total int64, isArray []bool) {
	if len(isArray) == 0 || !isArray[0] {
		if data == nil {
			data = struct{}{}
		}
		c.JSON(code, &Response{Code: errCode, Message: msg, Data: data})
	} else {
		if data == nil {
			data = []struct{}{}
		}
		c.JSON(code, &ArrayResponse{Code: errCode, Message: msg, Total: total, Data: data})
	}

	c.Abort()
}

func errorResponse(ctx *gin.Context, code int, errCode int, err error, msg string, isArray []bool) {
	if err != nil {
		msg += ": " + err.Error()
	}
	response(ctx, code, errCode, msg, nil, 0, isArray)
	msg = ctx.Request.Method + " " + ctx.Request.URL.Path + ": " + msg
	if err != nil {
		global.LOG.Error(msg + ": " + err.Error())
	} else {
		global.LOG.Error(msg)
	}
}

func success(ctx *gin.Context, msg string, data any, total int64, isArray ...bool) {
	response(ctx, http.StatusOK, CodeOk, msg, data, total, isArray)
}

func Success(c *gin.Context, isArray ...bool) {
	success(c, SuccessMsg, struct{}{}, 0, isArray...)
}
func SuccessWithDataAndTotal(c *gin.Context, data any, total int64) {
	success(c, SuccessMsg, data, total, true)
}
func SuccessWithMessage(c *gin.Context, msg string) {
	success(c, msg, struct{}{}, 0)
}
func SuccessWithMessageAndData(c *gin.Context, msg string, data any) {
	success(c, msg, data, 0)
}

const IsCodeAndMessage = false

func SuccessWithData(c *gin.Context, data any) {
	if IsCodeAndMessage {
		success(c, SuccessMsg, data, 0, true)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
func SuccessWithArrayData(c *gin.Context, data any, total int64) {
	if IsCodeAndMessage {
		success(c, SuccessMsg, data, total, true)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
func SuccessAuth(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

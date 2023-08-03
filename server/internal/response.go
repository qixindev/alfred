package internal

import (
	"accounts/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Total     int64  `json:"total"`
	Data      any    `json:"data"`
	PageTotal int64  `json:"pageTotal"`
}

func ErrorSqlResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    struct{}{},
	})
}

func ErrorNotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: msg,
		Data:    struct{}{},
	})
}

func ErrReqPara(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "req para err: failed to bind json, " + err.Error(),
		Data:    struct{}{},
	})
	global.LOG.Error("req para err: " + err.Error())
}

func ErrReqParaWithMsg(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: msg,
		Data:    struct{}{},
	})
	global.LOG.Error("req para err: " + err.Error())
}

func ErrReqParaCustom(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "req para err: " + err,
		Data:    struct{}{},
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    struct{}{},
	})
}

func SuccessWithDataAndTotal(c *gin.Context, data any, pageTotal int64, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Message:   "操作成功",
		Data:      data,
		Total:     total,
		PageTotal: pageTotal,
	})
}

func SuccessWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    map[string]interface{}{},
	})
}

func SuccessWithMessageAndData(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

package internal

import (
	"accounts/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Total   uint   `json:"total"`
	Data    any    `json:"data"`
}

func ErrorSqlResponse(msg string) Response {
	return Response{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    struct{}{},
	}
}

func ErrReqPara(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "req para err",
		Data:    struct{}{},
	})
	global.LOG.Error("req para err: " + err.Error())
}

func ErrReqParaCustom(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "req para err",
		Data:    struct{}{},
	})
	global.LOG.Error("req para err: " + err)
}

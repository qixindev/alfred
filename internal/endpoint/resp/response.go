package resp

import (
	"accounts/pkg/global"
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

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    data,
	})
}

func SuccessWithDataAndTotal(c *gin.Context, data any, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
		Total:   total,
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
func ErrorUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Data:    struct{}{},
	})
}

func ErrorForbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: msg,
		Data:    struct{}{},
	})
}

func ErrorInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    struct{}{},
	})
}

func ErrorServiceUnavailable(c *gin.Context, msg string) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Code:    http.StatusServiceUnavailable,
		Message: msg,
		Data:    struct{}{},
	})
}

func SuccessWithCustomCode(c *gin.Context, code int, msg string, data any) {
	c.JSON(code, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

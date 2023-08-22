package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 请求相关错误

func ErrorRequestWithMsg(c *gin.Context, msg string, isArray ...bool) {
	errorResponse(c, http.StatusBadRequest, CodeRequest, nil, msg, isArray)
}
func ErrorRequest(c *gin.Context, err error, isArray ...bool) {
	errorResponse(c, http.StatusBadRequest, CodeRequest, err, "req para err", isArray)
}

func ErrorUnauthorized(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusUnauthorized, CodeUnauthorized, err, msg, isArray)
}
func ErrorIamPermissionDeny(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusForbidden, CodeIamDeny, err, msg, isArray)
}
func ErrorForbidden(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusForbidden, CodeForbidden, err, msg, isArray)
}
func ErrorNotFound(c *gin.Context, msg string, isArray ...bool) {
	errorResponse(c, http.StatusNotFound, CodeNotFound, nil, msg, isArray)
}

func ErrReqPara(c *gin.Context, err error, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, http.StatusInternalServerError,
		nil, "req para err: "+err.Error(), isArray)
}

func ErrReqParaCustom(c *gin.Context, err string, isArray ...bool) {
	errorResponse(c, http.StatusBadRequest, http.StatusBadRequest, nil, "req para err: "+err, isArray)
}

func ErrorNotLogin(c *gin.Context, err error) {
	errorResponse(c, http.StatusUnauthorized, CodeNotLogin, err, "user not login", []bool{})
}

package resp

import (
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const (
	CodeOk          = 200
	SuccessMsg      = ""
	CodeUnknown     = 1000
	CodeSaveSession = 1001

	CodeSqlFirst  = 2000
	CodeSqlSelect = 2001
	CodeSqlCreate = 2002
	CodeSqlUpdate = 2003
	CodeSqlDelete = 2004

	CodeRequest      = 3000 // 请求参数错误
	CodeUnauthorized = 3001 // 未授权
	CodeIamDeny      = 3002 // 无权访问
)

func ErrorUnknown(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeUnknown, err, msg, isArray)
}

func ErrorSaveSession(c *gin.Context, err error, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSaveSession, err, "save session err", isArray)
}

// sql相关错误

func ErrorSqlFirst(c *gin.Context, err error, msg string, isArray ...bool) {
	if err == gorm.ErrRecordNotFound {
		errorResponse(c, http.StatusNotFound, CodeSqlFirst, err, msg, isArray)
	} else {
		errorResponse(c, http.StatusInternalServerError, CodeSqlFirst, err, msg, isArray)
	}
}
func ErrorSqlSelect(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSqlSelect, err, msg, isArray)
}
func ErrorSqlCreate(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSqlCreate, err, msg, isArray)
}
func ErrorSqlUpdate(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSqlUpdate, err, msg, isArray)
}
func ErrorSqlDelete(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSqlDelete, err, msg, isArray)
}

// 请求相关错误

func ErrorRequest(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusBadRequest, CodeRequest, err, msg, isArray)
}
func ErrorUnauthorized(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusUnauthorized, CodeUnauthorized, err, msg, isArray)
}
func ErrorIamPermissionDeny(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusForbidden, CodeIamDeny, err, msg, isArray)
}
func ErrorNotFound(c *gin.Context, msg string, isArray ...bool) {
	errorResponse(c, http.StatusNotFound, CodeIamDeny, nil, msg, isArray)
}

func ErrorSqlResponse(c *gin.Context, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, http.StatusInternalServerError, nil, msg, isArray)
}

func ErrReqPara(c *gin.Context, err error, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, http.StatusInternalServerError,
		nil, "req para err: failed to bind json, "+err.Error(), isArray)
	global.LOG.Error("req para err: " + err.Error())
}

func ErrReqParaCustom(c *gin.Context, err string, isArray ...bool) {
	errorResponse(c, http.StatusBadRequest, http.StatusBadRequest, nil, "req para err: "+err, isArray)
}

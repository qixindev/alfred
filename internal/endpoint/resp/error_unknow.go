package resp

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

const (
	CodeOk          = 200
	SuccessMsg      = ""
	CodeNotLogin    = 1000
	CodeSaveSession = 1001
	CodeNotFound    = 1002
	CodeUnknown     = 1003

	CodeSqlFirst  = 2000
	CodeSqlSelect = 2001
	CodeSqlCreate = 2002
	CodeSqlUpdate = 2003
	CodeSqlDelete = 2004

	CodeRequest      = 3000 // 请求参数错误
	CodeUnauthorized = 3001 // 未授权
	CodeForbidden    = 3002 // 无权访问
	CodeIamDeny      = 3003 // 无iam权限
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
	if err != nil && strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		errorResponse(c, http.StatusConflict, CodeSqlCreate, err, msg, isArray)
	} else {
		errorResponse(c, http.StatusInternalServerError, CodeSqlCreate, err, msg, isArray)
	}
}
func ErrorSqlUpdate(c *gin.Context, err error, msg string, isArray ...bool) {
	if err != nil && strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		errorResponse(c, http.StatusConflict, CodeSqlUpdate, err, msg, isArray)
	} else {
		errorResponse(c, http.StatusInternalServerError, CodeSqlUpdate, err, msg, isArray)
	}
}
func ErrorSqlDelete(c *gin.Context, err error, msg string, isArray ...bool) {
	if err == gorm.ErrForeignKeyViolated { // 外键依赖导致无法删除
		errorResponse(c, http.StatusConflict, CodeSqlDelete, err, msg, isArray)
	} else {
		errorResponse(c, http.StatusInternalServerError, CodeSqlDelete, err, msg, isArray)
	}
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
func ErrorForbidden(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusForbidden, CodeIamDeny, err, msg, isArray)
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

func ErrorNotLogin(c *gin.Context) {
	errorResponse(c, http.StatusInternalServerError, CodeNotLogin, nil, "user not login", []bool{})
}

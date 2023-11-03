package resp

import (
	"errors"
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
	CodeUpload      = 1004

	CodeSqlFirst  = 2000
	CodeSqlSelect = 2001
	CodeSqlCreate = 2002
	CodeSqlUpdate = 2003
	CodeSqlDelete = 2004

	CodeRequest      = 3000 // 请求参数错误
	CodeUnauthorized = 3001 // 未授权
	CodeForbidden    = 3002 // 无权访问
	CodeIamDeny      = 3003 // 无iam权限
	CodePassword     = 3004 // 密码错误
	CodeConflict     = 3005 // 资源冲突
	CodeValidate     = 3006 // 自定义请求参数错误
)

func ErrorUnknown(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeUnknown, err, msg, isArray)
}
func ErrorUpload(c *gin.Context, err error, msg string, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeUpload, err, msg, isArray)
}

func ErrorSaveSession(c *gin.Context, err error, isArray ...bool) {
	errorResponse(c, http.StatusInternalServerError, CodeSaveSession, err, "save session err", isArray)
}

// sql相关错误

func ErrorSqlFirst(c *gin.Context, err error, msg string, isArray ...bool) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
	if errors.Is(err, gorm.ErrForeignKeyViolated) { // 外键依赖导致无法删除
		errorResponse(c, http.StatusConflict, CodeSqlDelete, err, msg, isArray)
	} else {
		errorResponse(c, http.StatusInternalServerError, CodeSqlDelete, err, msg, isArray)
	}
}

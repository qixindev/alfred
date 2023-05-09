package internal

import (
	"accounts/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrReqPara(c *gin.Context, err error) {
	c.String(http.StatusBadRequest, "req para err")
	global.LOG.Error("req para err: " + err.Error())
}

func ErrReqParaCustom(c *gin.Context, err string) {
	c.String(http.StatusBadRequest, "req para err")
	global.LOG.Error("req para err: " + err)
}

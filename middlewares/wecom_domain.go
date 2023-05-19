package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WecomDomainCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/WW_verify_lPcI9c2g5B1tINHE.txt" {
			c.String(http.StatusOK, "lPcI9c2g5B1tINHE")
			c.Abort()
			return
		}
	}
}

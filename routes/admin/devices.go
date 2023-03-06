package admin

import "github.com/gin-gonic/gin"

func addAdminDevicesRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/:tenant/devices", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/devices/:deviceId", func(c *gin.Context) {

	})

	rg.POST("/admin/:tenant/devices", func(c *gin.Context) {

	})

	rg.PUT("/admin/:tenant/devices/:deviceId", func(c *gin.Context) {

	})

	rg.DELETE("/admin/:tenant/devices/:deviceId", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/devices/:deviceId/groups", func(c *gin.Context) {

	})

	rg.DELETE("/admin/:tenant/devices/:deviceId/groups/:groupId", func(c *gin.Context) {

	})
}

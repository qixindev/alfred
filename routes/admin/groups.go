package admin

import "github.com/gin-gonic/gin"

func addAdminGroupsRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/:tenant/groups", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/groups/:groupId", func(c *gin.Context) {

	})

	rg.POST("/admin/:tenant/groups", func(c *gin.Context) {

	})

	rg.PUT("/admin/:tenant/groups/:groupId", func(c *gin.Context) {

	})

	rg.DELETE("/admin/:tenant/groups/:groupId", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/groups/:groupId/members", func(c *gin.Context) {

	})
}

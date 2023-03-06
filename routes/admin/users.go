package admin

import "github.com/gin-gonic/gin"

func addAdminUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/admin/:tenant/users", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/users/:userId", func(c *gin.Context) {

	})

	rg.POST("/admin/:tenant/users", func(c *gin.Context) {

	})

	rg.PUT("/admin/:tenant/users/:userId", func(c *gin.Context) {

	})

	rg.DELETE("/admin/:tenant/users/:userId", func(c *gin.Context) {

	})

	rg.GET("/admin/:tenant/users/:userId/groups", func(c *gin.Context) {

	})

	rg.DELETE("/admin/:tenant/users/:userId/groups/:groupId", func(c *gin.Context) {

	})
}

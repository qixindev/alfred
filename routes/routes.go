package routes

import "github.com/gin-gonic/gin"

func AddRoutes(rg *gin.RouterGroup) {
	tenantRoutes := rg.Group("/:tenant", MultiTenancy)
	addLoginRoutes(tenantRoutes)
	addUsersRoutes(tenantRoutes)
	addOAuth2Routes(tenantRoutes)
}

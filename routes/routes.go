package routes

import (
	"accounts/middlewares"
	"accounts/routes/admin"
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	tenantRoutes := rg.Group("/:tenant", middlewares.MultiTenancy)
	addLoginRoutes(tenantRoutes)
	addUsersRoutes(tenantRoutes)
	addOAuth2Routes(tenantRoutes)

	admin.AddAdminRoutes(rg)
}

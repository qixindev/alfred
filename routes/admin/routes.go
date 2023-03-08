package admin

import (
	"accounts/middlewares"
	"github.com/gin-gonic/gin"
)

func AddAdminRoutes(rg *gin.RouterGroup) {
	tenantRoutes := rg.Group("/admin/:tenant", middlewares.MultiTenancy, middlewares.AuthorizedAdmin)
	addAdminGroupsRoutes(tenantRoutes)
	addAdminUsersRoutes(tenantRoutes)
	addAdminDevicesRoutes(tenantRoutes)
	addAdminProvidersRoutes(tenantRoutes)
	addAdminClientsRoutes(tenantRoutes)

	adminRoutes := rg.Group("/admin")
	addAdminTenantsRoutes(adminRoutes)
}

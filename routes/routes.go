package routes

import (
	"accounts/middlewares"
	"accounts/routes/admin"
	"accounts/routes/iam"
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	tenantRoutes := rg.Group("/:tenant", middlewares.MultiTenancy)
	addLoginRoutes(tenantRoutes)
	addUsersRoutes(tenantRoutes)
	addOAuth2Routes(tenantRoutes)

	iamRoutes := rg.Group("/:tenant/iam/clients/:client", middlewares.MultiTenancy)
	iam.AddIamRoutes(iamRoutes)

	admin.AddAdminRoutes(rg)
}

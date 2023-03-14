package routes

import (
	_ "accounts/docs"
	"accounts/middlewares"
	"accounts/routes/admin"
	"accounts/routes/iam"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tenantRoutes := rg.Group("/:tenant", middlewares.MultiTenancy)
	addLoginRoutes(tenantRoutes)
	addUsersRoutes(tenantRoutes)
	addOAuth2Routes(tenantRoutes)

	iamRoutes := rg.Group("/:tenant/iam/clients/:client", middlewares.MultiTenancy)
	iam.AddIamRoutes(iamRoutes)

	admin.AddAdminRoutes(rg)
}

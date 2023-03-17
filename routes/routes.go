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
	api := rg.Group("/accounts")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		tenantRoutes := api.Group("/:tenant", middlewares.MultiTenancy)
		addLoginRoutes(tenantRoutes)
		addUsersRoutes(tenantRoutes)
		addOAuth2Routes(tenantRoutes)

		iamRoutes := api.Group("/:tenant/iam/clients/:client", middlewares.MultiTenancy)
		iam.AddIamRoutes(iamRoutes)
		admin.AddAdminRoutes(api)
	}
}

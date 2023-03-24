package router

import (
	_ "accounts/docs"
	"accounts/middlewares"
	"accounts/router/admin"
	"accounts/router/iam"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AddRoutes(r *gin.Engine) {
	AddWebRoutes(r)
	r.GET("/accounts/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tenantApi := r.RouterGroup.Group("/accounts/:tenant", middlewares.MultiTenancy)
	{
		addLoginRoutes(tenantApi)
		addUsersRoutes(tenantApi)
		addOAuth2Routes(tenantApi)
	}

	adminApi := r.RouterGroup.Group("/accounts/admin/:tenant", middlewares.MultiTenancy, middlewares.AuthorizedAdmin)
	{
		admin.AddAdminGroupsRoutes(adminApi)
		admin.AddAdminUsersRoutes(adminApi)
		admin.AddAdminDevicesRoutes(adminApi)
		admin.AddAdminProvidersRoutes(adminApi)
		admin.AddAdminClientsRoutes(adminApi)
	}

	admin.AddAdminTenantsRoutes(r.RouterGroup.Group("/accounts/admin")) // all tenants
	iam.AddIamRoutes(r.RouterGroup.Group("/accounts/:tenant/iam/clients/:client", middlewares.MultiTenancy))
}

func AddWebRoutes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./web/.output/public", false)))
	r.StaticFile("/", "./web/.output/public/index.html")
}

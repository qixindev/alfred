package internal

import (
	_ "alfred/docs"
	"alfred/internal/controller"
	"alfred/internal/controller/admin"
	"alfred/internal/controller/authentication"
	"alfred/internal/controller/iam"
	"alfred/internal/controller/reset"
	"alfred/internal/endpoint/resp"
	"alfred/pkg/middlewares"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
)

func AddRoutes(r *gin.Engine) {
	AddWebRoutes(r)
	r.Use(middlewares.AccessJsMiddleware())
	r.Use(middlewares.WecomDomainCheck())
	r.GET("/accounts/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tenantApi := r.RouterGroup.Group("/accounts/:tenant", middlewares.MultiTenancy)
	{
		authentication.AddLoginRoutes(tenantApi)
		authentication.AddUsersRoutes(tenantApi)
		controller.AddOAuth2Routes(tenantApi)
		controller.AddMsgRouter(tenantApi)
		reset.AddResetRouter(tenantApi)
	}

	r.GET("/accounts/login/providers/callback", authentication.ProviderCallback) // 验证第三方登录是否成功
	adminApi := r.RouterGroup.Group("/accounts/admin/:tenant", middlewares.MultiTenancy, middlewares.AuthorizedAdmin)
	{
		admin.AddAdminGroupsRoutes(adminApi)
		admin.AddAdminUsersRoutes(adminApi)
		admin.AddAdminDevicesRoutes(adminApi)
		admin.AddAdminProvidersRoutes(adminApi)
		admin.AddAdminSmsRoutes(adminApi)
		admin.AddAdminClientsRoutes(adminApi)
		admin.AddClientUserRoute(adminApi)
	}

	adminRouter := r.RouterGroup.Group("/accounts/admin", middlewares.MultiTenancy)
	admin.AddAdminTenantsRoutes(adminRouter) // all tenants

	iamRouter := r.RouterGroup.Group("/accounts/:tenant/iam/clients/:client", middlewares.MultiTenancy, middlewares.AuthorizedAdmin)
	iam.AddIamRoutes(iamRouter)
}

func AddWebRoutes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./web/.output/public", false)))
	r.GET("/dashboard/*any", func(c *gin.Context) {
		c.File("./web/.output/public/index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/accounts") {
			resp.ErrorNotFound(c, "no such router")
			return
		}
		c.File("./web/.output/public/404.html")
	})
}

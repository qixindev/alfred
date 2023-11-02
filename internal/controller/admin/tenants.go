package admin

import (
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/global"
	"alfred/pkg/middlewares"
	"alfred/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ListTenants godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	list tenants
//	@Tags			admin-tenants
//	@Success		200
//	@Router			/accounts/admin/tenants [get]
func ListTenants(c *gin.Context) {
	var tenants []model.Tenant
	username := sessions.Default(c).Get("user")
	if err := global.DB.Model(model.User{}).Select("t.id, t.name, users.role").
		Joins("LEFT JOIN tenants as t ON t.id = users.tenant_id").
		Where("users.username = ?", username).
		Find(&tenants).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list tenants err", true)
		return
	}

	res := make([]model.Tenant, 0)
	for _, tenant := range tenants {
		if tenant.Role == "admin" || tenant.Role == "owner" {
			res = append(res, tenant)
		}
	}
	resp.SuccessWithArrayData(c, utils.Filter(res, model.Tenant2Dto), 0)
}

func ListAllTenants(c *gin.Context) {
	var tenants []model.Tenant
	if err := global.DB.Find(&tenants).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list all tenant err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(tenants, model.Tenant2Dto), 0)
}

// GetTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	get tenants
//	@Tags			admin-tenants
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [get]
func GetTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	resp.SuccessWithData(c, tenant.Dto())
}

// NewTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	new tenants
//	@Tags			admin-tenants
//	@Success		200
//	@Router			/accounts/admin/tenants [post]
func NewTenant(c *gin.Context) {
	var tenant model.Tenant
	if err := c.BindJSON(&tenant); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if tenant.Sub == "" {
		resp.ErrReqParaCustom(c, "sub should not be null")
		return
	}

	if err := global.DB.Create(&tenant).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "new tenant err")
		return
	}

	if err := service.CopyUser(tenant.Sub, tenant.Id); err != nil {
		resp.ErrorSqlCreate(c, err, "copy tenant err")
		return
	}

	resp.SuccessWithData(c, tenant.Dto())
}

// UpdateTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	update tenants
//	@Tags			admin-tenants
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [put]
func UpdateTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	var t model.Tenant
	if err := c.BindJSON(&t); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	tenant.Name = t.Name
	if err := global.DB.Save(&tenant).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update tenant err")
		return
	}
	resp.SuccessWithData(c, tenant.Dto())
}

// DeleteTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [delete]
func DeleteTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	if err := service.DeleteTenant(tenant); err != nil {
		resp.ErrorSqlDelete(c, err, "delete tenant err")
		return
	}
	resp.Success(c)
}

// DeleteTenantSecret godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenantId		path	integer	true	"tenant"
//	@Param			secretId		path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId}/secrets/{secretId} [delete]
func DeleteTenantSecret(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}

	if err := utils.SetJWKS(tenant.Name, c.Param("secretId"), nil); err != nil {
		resp.ErrorUnknown(c, err, "delete jwks secret err")
	}

	resp.Success(c)
}

// NewTenantSecret godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId}/secrets [post]
func NewTenantSecret(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}

	var in struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	if err := utils.SetJWKS(tenant.Name, c.Param("secretId"), []byte(in.Secret)); err != nil {
		resp.ErrorUnknown(c, err, "create tenant secrete err")
	}

	resp.Success(c)
}

func AddAdminTenantsRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", middlewares.AuthorizedAdmin, ListTenants)
	rg.GET("/tenants/all", middlewares.AuthAccessToken, ListAllTenants)
	rg.GET("/tenants/:tenantId", middlewares.AuthorizedAdmin, GetTenant)
	rg.POST("/tenants", middlewares.AuthorizedAdmin, NewTenant)
	rg.PUT("/tenants/:tenantId", middlewares.AuthorizedAdmin, UpdateTenant)
	rg.DELETE("/tenants/:tenantId", middlewares.AuthorizedAdmin, DeleteTenant)
	rg.DELETE("/tenants/:tenantId/secrets/:secretId", middlewares.AuthorizedAdmin, DeleteTenantSecret)
	rg.POST("/tenants/:tenantId/secrets", middlewares.AuthorizedAdmin, NewTenantSecret)

	rg.GET("/:tenant/clients/:clientId/page/login", GetLoginPage)                                 // 获取登录页面配置
	rg.PUT("/:tenant/clients/:clientId/page/login", middlewares.AuthorizedAdmin, UpdateLoginPage) // 更新登录页面配置
}

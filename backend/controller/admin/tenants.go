package admin

import (
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/middlewares"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ListTenants
// @Summary	list tenants
// @Tags	admin-tenants
// @Success	200
// @Router	/accounts/admin/tenants [get]
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

// GetTenant
// @Summary	get tenants
// @Tags	admin-tenants
// @Param	tenantId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/tenants/{tenantId} [get]
func GetTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	resp.SuccessWithData(c, tenant.Dto())
}

// NewTenant
// @Summary	new tenants
// @Tags	admin-tenants
// @Success	200
// @Router	/accounts/admin/tenants [post]
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

	if _, err := utils.LoadRsaPublicKeys(tenant.Name); err != nil {
		resp.ErrorUnknown(c, err, "create jwk error")
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

// UpdateTenant
// @Summary	update tenants
// @Tags	admin-tenants
// @Param	tenantId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/tenants/{tenantId} [put]
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

// DeleteTenant
// @Summary	delete tenants
// @Tags	admin-tenants
// @Param	tenantId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/tenants/{tenantId} [delete]
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

// DeleteTenantSecret
// @Summary	delete tenants
// @Tags	admin-tenants
// @Param	tenantId	path	integer	true	"tenant"
// @Param	secretId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/tenants/{tenantId}/secrets/{secretId} [delete]
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

// NewTenantSecret
// @Summary	delete tenants
// @Tags	admin-tenants
// @Param	tenantId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/tenants/{tenantId}/secrets [post]
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
	adminTenant := rg.Group("", middlewares.AuthorizedAdmin)
	{
		adminTenant.GET("/tenants", ListTenants)
		adminTenant.GET("/tenants/:tenantId", GetTenant)
		adminTenant.POST("/tenants", NewTenant)
		adminTenant.PUT("/tenants/:tenantId", UpdateTenant)
		adminTenant.DELETE("/tenants/:tenantId", DeleteTenant)
		adminTenant.DELETE("/tenants/:tenantId/secrets/:secretId", DeleteTenantSecret)
		adminTenant.POST("/tenants/:tenantId/secrets", NewTenantSecret)

		adminTenant.PUT("/:tenant/page/login", UpdateLoginPage)               // 更新登录页面配置
		adminTenant.PUT("/:tenant/proto", UpdateTenantProto)                  // 更新用户隐私协议
		adminTenant.PUT("/:tenant/picture/:type/upload", UploadTenantPicture) // 更新图片
	}

	rg.GET("/tenants/all", middlewares.AuthAccessToken, ListAllTenants)
	rg.GET("/:tenant/page/login", GetLoginPage) // 获取登录页面配置
	rg.GET("/:tenant/proto", GetTenantProto)    // 获取用户隐私协议
}

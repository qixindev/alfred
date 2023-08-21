package admin

import (
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
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

// ListUserTenants godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	list tenants
//	@Tags			admin-tenants
//	@Param			user	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/users/{user} [get]
func ListUserTenants(c *gin.Context) {
	userId := c.Param("user")
	var tenants []model.Tenant
	if err := global.DB.Where("sub = ?", userId).Find(&tenants).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list tenant users err", true)
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
	if err := global.DB.Delete(&tenant).Error; err != nil {
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
	rg.GET("/tenants", ListTenants)
	rg.GET("/tenants/users/:user", ListUserTenants)
	rg.GET("/tenants/:tenantId", GetTenant)
	rg.POST("/tenants", NewTenant)
	rg.PUT("/tenants/:tenantId", UpdateTenant)
	rg.DELETE("/tenants/:tenantId", DeleteTenant)
	rg.DELETE("/tenants/:tenantId/secrets/:secretId", DeleteTenantSecret)
	rg.POST("/tenants/:tenantId/secrets", NewTenantSecret)
}

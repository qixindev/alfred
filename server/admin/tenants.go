package admin

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListTenants godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	list tenants
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants [get]
func ListTenants(c *gin.Context) {
	var tenants []models.Tenant
	if err := global.DB.Find(&tenants).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get tenants err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(tenants, models.Tenant2Dto))
}

// GetTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	get tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [get]
func GetTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get tenant err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, tenant.Dto())
}

// NewTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	new tenants
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants [post]
func NewTenant(c *gin.Context) {
	var tenant models.Tenant
	if err := c.BindJSON(&tenant); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	if err := global.DB.Create(&tenant).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new tenants err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, tenant.Dto())
}

// UpdateTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	update tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [put]
func UpdateTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get tenant err: " + err.Error())
		return
	}
	var t models.Tenant
	if err := c.BindJSON(&t); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	tenant.Name = t.Name
	if err := global.DB.Save(&tenant).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update tenant err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, tenant.Dto())
}

// DeleteTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId} [delete]
func DeleteTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get tenant err: " + err.Error())
		return
	}
	if err := global.DB.Delete(&tenant).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete tenant err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteTenantSecret godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId}/secrets/{secretId} [delete]
func DeleteTenantSecret(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get tenant err: " + err.Error())
		return
	}

	if err := utils.SetJWKS(tenant.Name, c.Param("secretId"), nil); err != nil {
		c.String(http.StatusInternalServerError, "delete failed")
		global.LOG.Error("delete tenant secret err: " + err.Error())
	}

	c.Status(http.StatusNoContent)
}

// NewTenantSecret godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenantId}/secrets [post]
func NewTenantSecret(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "id = ?", tenantId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get tenant err: " + err.Error())
		return
	}

	var in struct {
		Secret string `json:"secret"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	if err := utils.SetJWKS(tenant.Name, c.Param("secretId"), []byte(in.Secret)); err != nil {
		c.String(http.StatusInternalServerError, "delete failed")
		global.LOG.Error("delete tenant secret err: " + err.Error())
	}

	c.Status(http.StatusNoContent)
}

func AddAdminTenantsRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", ListTenants)
	rg.GET("/tenants/:tenantId", GetTenant)
	rg.POST("/tenants", NewTenant)
	rg.PUT("/tenants/:tenantId", UpdateTenant)
	rg.DELETE("/tenants/:tenantId", DeleteTenant)
	rg.DELETE("/tenants/:tenantId/secrets/:secretId", DeleteTenantSecret)
	rg.POST("/tenants/:tenantId/secrets", NewTenantSecret)
}

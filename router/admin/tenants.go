package admin

import (
	"accounts/global"
	"accounts/models"
	"accounts/router/internal"
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
	if global.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
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
	if global.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
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
	if global.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if err := global.DB.Delete(&tenant).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete tenant err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func addAdminTenantsRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", ListTenants)
	rg.GET("/tenants/:tenantId", GetTenant)
	rg.POST("/tenants", NewTenant)
	rg.PUT("/tenants/:tenantId", UpdateTenant)
	rg.DELETE("/tenants/:tenantId", DeleteTenant)
}

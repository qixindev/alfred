package admin

import (
	"accounts/data"
	"accounts/models"
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
	if data.DB.Find(&tenants).Error != nil {
		c.Status(http.StatusInternalServerError)
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
	if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
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
	err := c.BindJSON(&tenant)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if data.DB.Create(&tenant).Error != nil {
		c.Status(http.StatusConflict)
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
	if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var t models.Tenant
	err := c.BindJSON(&t)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	tenant.Name = t.Name
	if data.DB.Save(&tenant).Error != nil {
		c.Status(http.StatusInternalServerError)
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
	if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if data.DB.Delete(&tenant).Error != nil {
		c.Status(http.StatusInternalServerError)
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

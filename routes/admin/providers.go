package admin

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListProviders godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	list provider
//	@Tags			provider
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []models.Provider
	if middlewares.TenantDB(c).Find(&providers).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(providers, models.Provider2Dto))
}

// GetProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	get provider
//	@Tags			provider
//	@Param			tenant		path	string	true	"tenant"
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/providers/{providerId} [get]
func GetProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, provider.Dto())
}

// NewProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	new provider
//	@Tags			provider
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/providers [post]
func NewProvider(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var provider models.Provider
	err := c.BindJSON(&provider)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	provider.TenantId = tenant.Id
	if data.DB.Create(&provider).Error != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, provider.Dto())
}

// UpdateProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	update provider
//	@Tags			provider
//	@Param			tenant		path	string	true	"tenant"
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/providers/{providerId} [put]
func UpdateProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var p models.Provider
	err := c.BindJSON(&p)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	provider.Name = p.Name
	if data.DB.Save(&provider).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, provider.Dto())
}

// DeleteProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	delete provider
//	@Tags			provider
//	@Param			tenant		path	string	true	"tenant"
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/providers/{providerId} [delete]
func DeleteProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if data.DB.Delete(&provider).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func addAdminProvidersRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:providerId", GetProvider)
	rg.POST("/providers", NewProvider)
	rg.PUT("/providers/:providerId", UpdateProvider)
	rg.DELETE("/providers/:providerId", DeleteProvider)
}

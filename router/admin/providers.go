package admin

import (
	"accounts/global"
	"accounts/middlewares"
	"accounts/models"
	"accounts/router/internal"
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
//	@Router			/accounts/admin/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []models.Provider
	if err := middlewares.TenantDB(c).Find(&providers).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get provider err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [get]
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
//	@Router			/accounts/admin/{tenant}/providers [post]
func NewProvider(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var provider models.Provider
	if err := c.BindJSON(&provider); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	provider.TenantId = tenant.Id
	if err := global.DB.Create(&provider).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new provider err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [put]
func UpdateProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var p models.Provider
	if err := c.BindJSON(&p); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	provider.Name = p.Name
	if err := global.DB.Save(&provider).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update provider err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [delete]
func DeleteProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "id = ?", providerId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if err := global.DB.Delete(&provider).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete provider err: " + err.Error())
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

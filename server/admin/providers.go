package admin

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
	"accounts/server/service"
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
	if err := internal.TenantDB(c).Find(&providers).Error; err != nil {
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
	tenant := internal.GetTenant(c)
	providerId := c.Param("providerId")
	var p models.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	res, err := service.GetProvider(tenant.Id, p.Id, p.Type)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get provider config err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// NewProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	new provider
//	@Tags			provider
//	@Param			tenant	path	string	true	"tenant"
//	@Param			req		body	object	true	"body"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers [post]
func NewProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var provider models.Provider
	if err := c.BindJSON(&provider); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	if !service.IsValidType(provider.Type) {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "invalid type"})
		return
	}
	provider.TenantId = tenant.Id
	if err := global.DB.Create(&provider).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new provider err: " + err.Error())
		return
	}

	if err := service.CreateProviderConfig(provider); err != nil {
		c.JSON(http.StatusInternalServerError, &gin.H{"message": "failed to create provider config"})
		global.LOG.Error("create provider config err: " + err.Error())
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
	var p models.Provider
	if err := c.BindJSON(&p); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	name := p.Name
	if !service.IsValidType(p.Type) {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "invalid type"})
		return
	}

	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	p.Name = name
	if err := global.DB.Save(&p).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update provider err: " + err.Error())
		return
	}

	if err := service.UpdateProviderConfig(p); err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update provider config err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, p.Dto())
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
	if err := internal.TenantDB(c).First(&provider, "id = ?", providerId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}
	if err := global.DB.Delete(&provider).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete provider err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func AddAdminProvidersRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:providerId", GetProvider)
	rg.POST("/providers", NewProvider)
	rg.PUT("/providers/:providerId", UpdateProvider)
	rg.DELETE("/providers/:providerId", DeleteProvider)
}

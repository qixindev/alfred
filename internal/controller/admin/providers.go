package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/req"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListProviders godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	list provider
//	@Tags			provider
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []model.Provider
	if err := internal.TenantDB(c).Find(&providers).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(providers, model.Provider2Dto))
}

// GetProvider godoc
//
//	@Summary	provider
//	@Schemes
//	@Description	get provider
//	@Tags			provider
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [get]
func GetProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerId := c.Param("providerId")
	var p model.Provider
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

// GetProviderUsers godoc
//
//	@Summary	get provider user list
//	@Schemes
//	@Description	get provider user list
//	@Tags			provider
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			providerId	path	integer	true	"provider id"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers/{providerId}/users [get]
func GetProviderUsers(c *gin.Context) {
	providerId := c.Param("providerId")
	tenant := internal.GetTenant(c)
	var p model.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	res, err := service.GetProviderUsers(tenant.Id, p.Id)
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
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			req		body	object	true	"body"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers [post]
func NewProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var provider req.ReqProvider
	if err := c.BindJSON(&provider); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	provider.Id = 0
	provider.TenantId = tenant.Id
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [put]
func UpdateProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerId := c.Param("providerId")
	var p req.ReqProvider
	if err := c.BindJSON(&p); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	p.TenantId = tenant.Id
	p.ProviderId = utils.StrToUint(providerId)
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			providerId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/providers/{providerId} [delete]
func DeleteProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider model.Provider
	if err := internal.TenantDB(c).First(&provider, "id = ?", providerId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	if err := service.DeleteProviderConfig(provider); err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete provider config err: " + err.Error())
		return
	}

	if err := internal.TenantDB(c).Where("id = ?", providerId).Delete(&provider).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete provider err: " + err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func AddAdminProvidersRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:providerId", GetProvider)
	rg.GET("/providers/:providerId/users", GetProviderUsers)
	rg.POST("/providers", NewProvider)
	rg.PUT("/providers/:providerId", UpdateProvider)
	rg.DELETE("/providers/:providerId", DeleteProvider)
}

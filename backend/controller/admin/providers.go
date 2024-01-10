package admin

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/req"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"github.com/gin-gonic/gin"
)

// ListProviders
// @Summary	list provider
// @Tags	provider
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200
// @Router	/accounts/admin/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []model.Provider
	if err := internal.TenantDB(c).Find(&providers).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list provider err")
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(providers, model.Provider2Dto), 0)
}

// GetProvider
// @Summary	get provider
// @Tags	provider
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	providerId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/providers/{providerId} [get]
func GetProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerId := c.Param("providerId")
	var p model.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	res, err := service.GetProvider(tenant.Id, p.Id, p.Type)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider config err")
		return
	}

	resp.SuccessWithData(c, res)
}

// ListProviderUsers
// @Summary	get provider user list
// @Tags	provider
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	providerId	path	integer	true	"provider id"
// @Success	200
// @Router	/accounts/admin/{tenant}/providers/{providerId}/users [get]
func ListProviderUsers(c *gin.Context) {
	providerId := c.Param("providerId")
	tenant := internal.GetTenant(c)
	var p model.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err: ", true)
		return
	}

	res, err := service.GetProviderUsers(tenant.Id, p.Id)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list provider users err")
		return
	}

	resp.SuccessWithArrayData(c, res, 0)
}

// NewProvider
// @Summary	new provider
// @Tags	provider
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	req		body	object	true	"body"
// @Success	200
// @Router	/accounts/admin/{tenant}/providers [post]
func NewProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var provider req.Provider
	if err := c.BindJSON(&provider); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	provider.Id = 0
	provider.TenantId = tenant.Id
	if err := service.CreateProviderConfig(provider); err != nil {
		resp.ErrorSqlCreate(c, err, "create provider config err")
		return
	}
	resp.SuccessWithData(c, provider.Dto())
}

// UpdateProvider
// @Summary	update provider
// @Tags	provider
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	providerId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/providers/{providerId} [put]
func UpdateProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerId := c.Param("providerId")
	var p req.Provider
	if err := c.BindJSON(&p); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	p.TenantId = tenant.Id
	p.ProviderId = utils.StrToUint(providerId)
	if err := service.UpdateProviderConfig(p); err != nil {
		resp.ErrorSqlUpdate(c, err, "update provider config err")
		return
	}

	resp.SuccessWithData(c, p.Dto())
}

// DeleteProvider
// @Summary	delete provider
// @Tags	provider
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	providerId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/providers/{providerId} [delete]
func DeleteProvider(c *gin.Context) {
	providerId := c.Param("providerId")
	var provider model.Provider
	if err := internal.TenantDB(c).First(&provider, "id = ?", providerId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	if err := service.DeleteProviderConfig(provider); err != nil {
		resp.ErrorSqlDelete(c, err, "delete provider config err")
		return
	}

	if err := internal.TenantDB(c).Where("id = ?", providerId).Delete(&provider).Error; err != nil {
		resp.ErrorSqlDelete(c, err, "delete provider err")
		return
	}

	resp.Success(c)
}

func AddAdminProvidersRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:providerId", GetProvider)
	rg.GET("/providers/:providerId/users", ListProviderUsers)
	rg.POST("/providers", NewProvider)
	rg.PUT("/providers/:providerId", UpdateProvider)
	rg.DELETE("/providers/:providerId", DeleteProvider)
}

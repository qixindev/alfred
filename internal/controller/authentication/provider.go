package authentication

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/internal/service/auth"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// ListProviders godoc
//
//	@Summary	List all providers
//	@Schemes
//	@Description	list login providers
//	@Tags			login
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Success		200		{object}	[]dto.ProviderDto
//	@Router			/accounts/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []model.Provider
	if err := internal.TenantDB(c).Find(&providers).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list provider err", true)
		return
	}
	resp.SuccessWithData(c, utils.Filter(providers, model.Provider2Dto))
}

// GetProvider godoc
//
//	@Summary	get a provider
//	@Schemes
//	@Description	get a login provider
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			provider	path		string	true	"provider"
//	@Success		200			{object}	dto.ProviderDto
//	@Router			/accounts/{tenant}/providers/{provider} [get]
func GetProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerName := c.Param("provider")
	authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	resp.SuccessWithData(c, authProvider.LoginConfig())
}

// LoginToProvider godoc
//
//	@Summary	login via a provider
//	@Schemes
//	@Description	login via a provider
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			provider	path	string	true	"provider"
//	@Param			phone		query	string	false	"phone"
//	@Param			next		query	string	false	"next"
//	@Success		302
//	@Router			/accounts/{tenant}/login/{provider} [get]
func LoginToProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerName := c.Param("provider")
	authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	authStr := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenant.Name, providerName)
	if providerName == "sms" {
		authStr = c.Query("phone")
	}

	location, err := authProvider.Auth(authStr, tenant.Id)
	if err != nil {
		resp.ErrorUnknown(c, err, "provider auth err")
		return
	}

	next := c.Query("next")
	if next != "" {
		session := sessions.Default(c)
		session.Set("next", next)
		_ = session.Save()
	}

	if providerName == "sms" {
		resp.Success(c)
		return
	}

	c.Redirect(http.StatusFound, location)
}

// ProviderCallback godoc
//
//	@Summary	provider callback
//	@Schemes
//	@Description	provider callback
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			provider	path	string	true	"provider"
//	@Param			code		query	string	true	"code"
//	@Param			phone		query	string	false	"phone"
//	@Param			next		query	string	false	"next"
//	@Success		302
//	@Success		200
//	@Router			/accounts/{tenant}/logged-in/{provider} [get]
func ProviderCallback(c *gin.Context) {
	providerName := c.Param("provider")
	var provider model.Provider
	if err := internal.TenantDB(c).First(&provider, "name = ?", providerName).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	authProvider, err := auth.GetAuthProvider(provider.TenantId, provider.Name)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get auth provider err")
		return
	}
	userInfo, err := authProvider.Login(c)
	if err != nil {
		resp.ErrorUnknown(c, err, "login err")
		return
	}

	var providerUser model.ProviderUser
	var user *model.User
	if err = internal.TenantDB(c).First(&providerUser, "provider_id = ? AND name = ?", provider.Id, userInfo.Sub).
		Error; errors.Is(err, gorm.ErrRecordNotFound) { // provider user不存在，直接创建
		user, err = service.BindLoginUser(userInfo, provider.TenantId, provider.Type)
		if err != nil {
			resp.ErrorSqlCreate(c, err, "bind login user err")
			return
		}
		providerUser.TenantId = provider.TenantId
		providerUser.ProviderId = provider.Id
		providerUser.UserId = user.Id
		providerUser.Name = userInfo.Sub
		if err = global.DB.Create(&providerUser).Error; err != nil {
			resp.ErrorSqlCreate(c, err, "create provider user err")
			return
		}
	} else if err != nil {
		resp.ErrorSqlSelect(c, err, "get user err")
		return
	} else {
		if err = global.DB.Where("id = ? AND tenant_id = ?", providerUser.UserId, provider.TenantId).
			First(&user).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "get user err")
			return
		}
	}

	session := sessions.Default(c)
	tenant := internal.GetTenant(c)
	session.Set("tenant", tenant.Name)
	session.Set("user", user.Username)
	next := utils.GetString(session.Get("next"))
	session.Delete("next")
	if err = session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
	}
	if next != "" {
		c.Redirect(http.StatusFound, next)
		return
	}
}

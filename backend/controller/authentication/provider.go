package authentication

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"alfred/backend/service/auth"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListProviders
// @Summary	List all providers
// @Tags	login
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200	{object}	[]dto.ProviderDto
// @Router	/accounts/{tenant}/providers [get]
func ListProviders(c *gin.Context) {
	var providers []model.Provider
	if err := internal.TenantDB(c).Find(&providers).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list provider err", true)
		return
	}
	resp.SuccessWithData(c, utils.Filter(providers, model.Provider2Dto))
}

// GetProvider
// @Summary	get a login provider
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	provider	path	string	true	"provider"
// @Success	200	{object}	dto.ProviderDto
// @Router	/accounts/{tenant}/providers/{provider} [get]
func GetProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerName := c.Param("provider")
	_, authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	resp.SuccessWithData(c, authProvider.LoginConfig())
}

type ProviderLogin struct {
	Redirect string `json:"redirect"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
	ClientId string `json:"clientId"`
	Tenant   string `json:"tenant"`
	TenantId uint   `json:"tenantId"`
}

// LoginToProvider
// @Summary	login via a provider
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	provider	path	string	true	"provider"
// @Param	phone		query	string	false	"phone"
// @Param	next		query	string	false	"next"
// @Param	callback	query	string	false	"callback url"
// @Success	302
// @Router	/accounts/{tenant}/providers/{provider}/login [get]
func LoginToProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerName := c.Param("provider")
	callbackUrl := c.Query("callback")
	provider, authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	state := uuid.NewString()
	authStr := utils.GetHostWithScheme(c) + "/redirect"
	if callbackUrl != "" {
		authStr = callbackUrl
	}
	if provider.Type == "sms" {
		state = c.Query("phone")
	}

	location, err := authProvider.Auth(authStr, state, tenant.Id)
	if err != nil {
		resp.ErrorUnknown(c, err, "provider auth err")
		return
	}

	if providerName == "sms" {
		resp.Success(c)
		return
	}

	loginInfo := ProviderLogin{
		Redirect: c.Query("next"),
		Provider: providerName,
		ClientId: "default",
		Tenant:   tenant.Name,
		TenantId: tenant.Id,
	}
	infoByte, err := json.Marshal(&loginInfo)
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to marshal provider info")
		return
	}

	if err = global.CodeCache.Set(state, infoByte); err != nil {
		resp.ErrorUnknown(c, err, "failed to set code")
		return
	}
	resp.SuccessWithData(c, &gin.H{"state": state, "location": location})
}

// ProviderCallback
// @Summary	provider callback
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	code		query	string	true	"code"
// @Param	state		query	string	false	"state"
// @Param	phone		query	string	false	"phone"
// @Success	200
// @Router	/accounts/login/providers/callback [get]
func ProviderCallback(c *gin.Context) {
	var provider model.Provider
	state := c.Query("state")
	loginInfo, err := global.CodeCache.Get(state)
	if err != nil {
		_ = global.CodeCache.Delete(state)
		resp.ErrorForbidden(c, err, "invalidate state")
		return
	}
	if err = global.CodeCache.Delete(state); err != nil {
		resp.ErrorUnknown(c, err, "failed to delete state")
		return
	}
	var stateInfo ProviderLogin
	if err = json.Unmarshal(loginInfo, &stateInfo); err != nil {
		resp.ErrorUnknown(c, err, "failed unmarshal cache login info")
		return
	}
	if err = global.DB.Where("tenant_id = ? AND name = ?", stateInfo.TenantId, stateInfo.Provider).
		First(&provider).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	_, authProvider, err := auth.GetAuthProvider(provider.TenantId, provider.Name)
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
	if err = global.DB.First(&providerUser, "provider_id = ? AND name = ? AND tenant_id = ?",
		provider.Id, userInfo.Sub, stateInfo.TenantId).
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
	session.Set("tenant", stateInfo.Tenant)
	session.Set("user", user.Username)
	session.Delete("next")
	if err = session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
	}

	stateInfo.Type = provider.Type
	resp.SuccessWithData(c, stateInfo)
}

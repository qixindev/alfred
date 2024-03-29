package authentication

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"alfred/backend/service/auth"
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
	authState := state
	authStr := utils.GetHostWithScheme(c) + "/redirect"
	if callbackUrl != "" {
		authStr = callbackUrl
	}
	if provider.Type == "sms" {
		authState = c.Query("phone")
	}

	location, err := authProvider.Auth(authStr, authState, tenant.Id)
	if err != nil {
		resp.ErrorUnknown(c, err, "provider auth err")
		return
	}

	loginInfo := global.StateInfo{
		State:     state,
		AuthState: authState,
		Type:      provider.Type,
		Provider:  providerName,
		Redirect:  c.Query("next"),
		ClientId:  "default",
		Tenant:    tenant.Name,
		TenantId:  tenant.Id,
	}
	if err = global.SetStateInfo(state, loginInfo); err != nil {
		resp.ErrorUnknown(c, err, "failed to set cache info")
		return
	}
	if provider.Type == "sms" {
		resp.SuccessWithData(c, &gin.H{"state": state})
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
	state := c.Query("state")
	stateInfo, err := global.GetAndDeleteStateInfo(state)
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to get or delete state")
		return
	}
	provider, authProvider, err := auth.GetAuthProvider(stateInfo.TenantId, stateInfo.Provider)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get auth provider err")
		return
	}
	userInfo, err := authProvider.Login(c.Query("code"), stateInfo)
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
	} else if err == nil {
		if err = global.DB.Where("id = ? AND tenant_id = ?", providerUser.UserId, provider.TenantId).
			First(&user).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "get user err")
			return
		}
	} else {
		resp.ErrorSqlSelect(c, err, "get user err")
		return
	}

	session := sessions.Default(c)
	session.Set("tenant", stateInfo.Tenant)
	session.Set("tenantId", stateInfo.TenantId)
	session.Set("client", stateInfo.ClientId)
	session.Set("user", user.Username)
	session.Set("userId", user.Id)
	session.Delete("next")
	if err = session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
	}

	resp.SuccessWithData(c, stateInfo)
}

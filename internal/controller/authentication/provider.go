package authentication

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/internal/service/auth"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	resp.SuccessWithData(c, authProvider.LoginConfig())
}

type ProviderLogin struct {
	Redirect string `json:"redirect"`
	Type     string `json:"type"`
	Client   string `json:"client"`
	Tenant   string `json:"tenant"`
	Location string `json:"location"`
}

// LoginToProvider
// @Summary	login via a provider
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	provider	path	string	true	"provider"
// @Param	phone		query	string	false	"phone"
// @Param	next		query	string	false	"next"
// @Success	302
// @Router	/accounts/{tenant}/login/{provider} [get]
func LoginToProvider(c *gin.Context) {
	tenant := internal.GetTenant(c)
	providerName := c.Param("provider")
	authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	authStr := utils.GetHostWithScheme(c) + "redirect"
	if providerName == "sms" {
		authStr = c.Query("phone")
	}

	location, err := authProvider.Auth(authStr, tenant.Id)
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
		Type:     providerName,
		Client:   "default",
		Tenant:   tenant.Name,
		Location: location,
	}
	infoByte, err := json.Marshal(&loginInfo)
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to marshal provider info")
		return
	}

	c.SetCookie("login-info", string(infoByte), 5*60, "/accounts", c.Request.Host, false, true)
	resp.SuccessWithData(c, &loginInfo)
}

// ProviderCallback
// @Summary	provider callback
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	provider	path	string	true	"provider"
// @Param	code		query	string	true	"code"
// @Param	phone		query	string	false	"phone"
// @Param	next		query	string	false	"next"
// @Success	302
// @Success	200
// @Router	/accounts/{tenant}/logged-in/{provider} [get]
func ProviderCallback(c *gin.Context) {
	providerName := c.Param("provider")
	var provider model.Provider
	if err := internal.TenantDB(c).First(&provider, "name = ?", providerName).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	loginInfo, err := c.Cookie("login-info")
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to get login info")
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
	session.Delete("next")
	if err = session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
	}

	c.SetCookie("login-info", "", -1, "/accounts", c.Request.Host, false, true)
	resp.SuccessWithData(c, loginInfo)
}

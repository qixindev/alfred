package authentication

import (
	"accounts/internal/controller/auth"
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/req"
	"accounts/internal/endpoint/resp"
	"accounts/pkg/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SmsRoute struct {
	internal.Api
}

func NewSmsRoute() *SmsRoute {
	return &SmsRoute{}
}

// LoginToSms godoc
//
//	@Summary	login via a provider
//	@Schemes
//	@Description	login via a provider
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"		default(default)
//	@Param			type		path	string	true	"connector type" default(tcloud)
//	@Param			next		query	string	false	"next"
//	@Success		302
//	@Router			/accounts/{tenant}/login/connectors/{type} [get]
func (s SmsRoute) LoginToSms(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typ := c.Param("type")
	authProvider, err := auth.GetAuthProvider(tenant.Id, typ)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenant.Name, typ)
	location, err := authProvider.Auth(redirectUri)
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
	c.Redirect(http.StatusFound, location)
}

// SmsCallback godoc
//
//	@Summary	login via a provider
//	@Schemes
//	@Description	login via a provider
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"		default(default)
//	@Param			type		path	string	true	"connector type" default(tcloud)
//	@Param			code		path	string	true	"verify code"
//	@Success		200
//	@Router			/accounts/{tenant}/login/connectors/{type}/verify [post]
func (s SmsRoute) SmsCallback(c *gin.Context) {
	var in req.SmsVerify
	if err := s.SetCtx(c).SetTenant().BindUriAndJson(&in).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	authProvider, err := auth.GetAuthProvider(s.Tenant.Id, in.Type)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), s.Tenant.Name, in.Type)
	location, err := authProvider.Auth(redirectUri)
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
	c.Redirect(http.StatusFound, location)
}

func AddSmsRoute(r *gin.RouterGroup) {
	s := NewSmsRoute()
	r.GET("/", s.LoginToSms)
	r.POST("/verify", s.LoginToSms)
}

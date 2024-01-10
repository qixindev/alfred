package controller

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/dto"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/middlewares"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAccessCode(c *gin.Context, client *model.Client) (string, error) {
	token, err := service.GetAccessToken(c, client)
	if err != nil {
		return "", err
	}
	code := uuid.NewString()
	tokenCode := model.TokenCode{
		Token:     token,
		Code:      code,
		CreatedAt: time.Now(),
		ClientId:  client.Id,
		TenantId:  client.TenantId,
	}
	if err = global.DB.Create(&tokenCode).Error; err != nil {
		return "", err
	}
	return code, nil
}

// GetAuthCode
// @Summary	oauth2 authorize
// @Tags	oauth2
// @Param	tenant	path	string	true	"tenant"
// @Param	client_id	query	string	true	"client_id"
// @Param	scope	query	string	true	"scope"
// @Param	response_type	query	string	true	"response_type"
// @Param	redirect_uri	query	string	true	"redirect_uri"
// @Param	state	query	string	false	"state"
// @Param	nonce	query	string	false	"nonce"
// @Success	302
// @Success	200
// @Router	/accounts/{tenant}/oauth2/auth [get]
func GetAuthCode(c *gin.Context) {
	// If logged in, 302 to redirect uri.
	// If not, return login form.
	clientId := c.Query("client_id")
	scope := c.Query("scope")
	responseType := c.Query("response_type")
	redirectUri := c.Query("redirect_uri")
	state := strings.TrimSpace(c.Query("state"))
	nonce := c.Query("nonce")
	tenant := internal.GetTenant(c)

	var client model.Client
	if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	if err := service.IsValidateUri(tenant.Id, client.Id, redirectUri); err != nil {
		resp.ErrorForbidden(c, err, "invalid redirect uri")
		return
	}

	if responseType == "code" {
		code, err := getAccessCode(c, &client)
		if err != nil {
			resp.ErrorSqlCreate(c, err, "create access code err")
			return
		}
		query := url.Values{}
		query.Add("code", code)
		if state != "" {
			query.Add("state", state)
		}
		location := fmt.Sprintf("%s?%s", redirectUri, query.Encode())
		c.Redirect(http.StatusFound, location)
		return
	} else if responseType == "token" {
		token, err := service.GetAccessToken(c, &client)
		if err != nil {
			resp.ErrorUnknown(c, err, "get access token err")
			return
		}
		query := url.Values{}
		query.Add("access_token", token)
		if state != "" {
			query.Add("state", state)
		}
		location := fmt.Sprintf("%s?%s", redirectUri, query.Encode())
		c.Redirect(http.StatusFound, location)
		return
	}

	resp.ErrorForbidden(c, nil, "Invalid response_type")
	fmt.Println(clientId, scope, responseType, redirectUri, state, nonce)
}

// GetDeviceCode
// @Summary	delete device groups
// @Tags	oauth2
// @Param	tenant		path	string	true	"tenant"
// @Param	client_id	query	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/oauth2/device/code [post]
func GetDeviceCode(c *gin.Context) {
	clientId := c.Query("client_id")
	scope := c.Query("scope")
	tenant := internal.GetTenant(c)
	var client model.Client
	if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}

	deviceCode := model.DeviceCode{
		TenantId: tenant.Id,
		Code:     uuid.NewString(),
		UserCode: utils.GetDeviceUserCode(),
		Status:   "verifying",
	}

	if err := internal.TenantDB(c).Create(&deviceCode).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "create device code err")
		return
	}

	verificationUri := utils.GetHostWithScheme(c) + "/accounts/admin/" + c.Param("tenant") + "/devices/code"
	resp.SuccessAuth(c, &gin.H{
		"deviceCode":              deviceCode.Code,
		"userCode":                deviceCode.UserCode,
		"clientId":                clientId,
		"expires_in":              2 * 60,
		"scope":                   scope,
		"verificationUri":         verificationUri,
		"verificationUriComplete": verificationUri + "/" + deviceCode.UserCode,
	})
}

// GetToken
// @Summary	oauth2 token
// @Tags	oauth2
// @Param	tenant			path	string	true	"tenant"
// @Param	client_id		query	string	true	"client_id"
// @Param	client_secret	query	string	false	"client_secret"
// @Param	code			query	string	false	"code"
// @Param	grant_type		query	string	true	"grant_type"
// @Param	redirect_uri	query	string	false	"redirect_uri"
// @Param	state			query	string	false	"state"
// @Param	nonce			query	string	false	"nonce"
// @Success	200	{object}	dto.AccessTokenDto
// @Router	/accounts/{tenant}/oauth2/token [get]
func GetToken(c *gin.Context) {
	clientId := c.Query("client_id")
	clientSecret := c.Query("client_secret")
	grantType := c.Query("grant_type")
	code := c.Query("code")
	redirectUri := c.Query("redirect_uri")
	state := strings.TrimSpace(c.Query("state"))
	nonce := c.Query("nonce")

	if grantType == "authorization_code" {
		tenant := internal.GetTenant(c)
		var client model.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid client id")
			return
		}
		var secret model.ClientSecret
		if err := internal.TenantDB(c).First(&secret, "client_id = ? AND secret = ?", client.Id, clientSecret).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid client secret")
			return
		}

		var tokenCode model.TokenCode
		if err := internal.TenantDB(c).First(&tokenCode, "code = ?", code).Error; err != nil {
			resp.ErrorForbidden(c, err, "invalid token code")
			return
		}
		service.ClearTokenCode(tokenCode.Code)
		accessToken := dto.AccessTokenDto{AccessToken: tokenCode.Token}
		resp.SuccessAuth(c, accessToken)
		return
	} else if grantType == "client_credential" {
		tenant := internal.GetTenant(c)
		var client model.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid client_id")
			return
		}
		var secret model.ClientSecret
		if err := internal.TenantDB(c).First(&secret, "client_id = ? AND secret = ?", client.Id, clientSecret).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid client secret")
			return
		}

		token, err := service.GetClientAccessToken(c, &client)
		if err != nil {
			resp.ErrorUnknown(c, err, "get accessToken err")
			return
		}
		resp.SuccessAuth(c, dto.AccessTokenDto{AccessToken: token})
		return
	} else if grantType == "urn:ietf:params:oauth:grant-type:device_code" {
		tenant := internal.GetTenant(c)
		var client model.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid client_id")
			return
		}

		var deviceCode model.DeviceCode
		if err := global.DB.First(&deviceCode, "tenant_id = ? AND user_code = ?", tenant.Id, code).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalid device code")
			return
		}
		if deviceCode.Status != "verified" {
			resp.ErrorForbidden(c, nil, "device code is unauthorized")
			return
		}

		token, err := service.GetClientAccessToken(c, &client)
		if err != nil {
			resp.ErrorUnknown(c, err, "generate accessToken err")
			return
		}

		service.ClearDeviceCode(deviceCode.UserCode)
		resp.SuccessAuth(c, dto.AccessTokenDto{AccessToken: token})
		return
	} else if grantType == "device_credential" {
		id := c.Query("device_id")
		secret := c.Query("device_secret")
		var device model.Device
		if err := internal.TenantDB(c).Where("id = ?", id).First(&device).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalidate device id")
			return
		}

		var deviceSecret model.DeviceSecret
		if err := internal.TenantDB(c).Where("device_id = ? AND secret = ?", id, secret).First(&deviceSecret).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "invalidate device secret")
			return
		}

		token, err := service.GetDeviceToken(c, &device)
		if err != nil {
			resp.ErrorUnknown(c, err, "generate accessToken err")
			return
		}

		resp.SuccessAuth(c, dto.AccessTokenDto{AccessToken: token})
		return
	}

	resp.ErrorForbidden(c, nil, "Invalid grant_type")
	fmt.Println(clientId, clientSecret, grantType, code, redirectUri, state, nonce)
}

// GetOpenidConfiguration
// @Summary	openid configuration
// @Tags	oauth2
// @Param	tenant	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/.well-known/openid-configuration [get]
func GetOpenidConfiguration(c *gin.Context) {
	tenant := internal.GetTenant(c)
	prefix := utils.GetHostWithScheme(c)
	conf := dto.OpenidConfigurationDto{
		Issuer:                            fmt.Sprintf("%s/%s", prefix, tenant.Name),
		AuthorizationEndpoint:             fmt.Sprintf("%s/%s/oauth2/auth", prefix, tenant.Name),
		TokenEndpoint:                     fmt.Sprintf("%s/%s/oauth2/token", prefix, tenant.Name),
		UserinfoEndpoint:                  fmt.Sprintf("%s/%s/me/profile", prefix, tenant.Name),
		JwksUri:                           fmt.Sprintf("%s/%s/.well-known/jwks.json", prefix, tenant.Name),
		ScopesSupported:                   []string{"openid", "profile", "email", "offline_access"},
		ResponseTypesSupported:            []string{"code", "id_token", "code id_token", "id_token token"},
		SubjectTypesSupported:             []string{"pairwise"},
		IdTokenSigningAlgValuesSupported:  []string{"RS256"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post"},
		ClaimsSupported:                   []string{"sub", "iss", "aud", "exp", "iat", "nonce", "name", "email"},
		RequestUriParameterSupported:      false,
	}
	resp.SuccessAuth(c, conf)
}

// GetJwks
// @Summary	jwk
// @Tags	oauth2
// @Param	tenant	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/.well-known/jwks.json [get]
func GetJwks(c *gin.Context) {
	tenant := internal.GetTenant(c)
	jwks, err := utils.LoadRsaPublicKeys(tenant.Name)
	if err != nil {
		resp.ErrorUnknown(c, err, "get jwks err")
		return
	}
	resp.SuccessAuth(c, jwks)
}

func AddOAuth2Routes(rg *gin.RouterGroup) {
	rg.GET("/oauth2/auth", middlewares.Authorized(true), GetAuthCode)
	rg.POST("/oauth2/device/code", GetDeviceCode)
	rg.GET("/oauth2/token", GetToken)
	rg.GET("/.well-known/openid-configuration", GetOpenidConfiguration)
	rg.GET("/.well-known/jwks.json", GetJwks)
}

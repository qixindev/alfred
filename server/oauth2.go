package server

import (
	"accounts/global"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"accounts/server/internal"
	"accounts/server/service"
	"accounts/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAccessCode(c *gin.Context, client *models.Client) (string, error) {
	token, err := internal.GetAccessToken(c, client)
	if err != nil {
		return "", err
	}
	code := uuid.NewString()
	tokenCode := models.TokenCode{
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

// GetAuthCode godoc
//
//	@Summary	oauth2 authorize
//	@Schemes
//	@Description	oauth2 authorize
//	@Tags			oauth2
//	@Param			tenant			path	string	true	"tenant"
//	@Param			client_id		query	string	true	"client_id"
//	@Param			scope			query	string	true	"scope"
//	@Param			response_type	query	string	true	"response_type"
//	@Param			redirect_uri	query	string	true	"redirect_uri"
//	@Param			state			query	string	false	"state"
//	@Param			nonce			query	string	false	"nonce"
//	@Success		302
//	@Success		200
//	@Router			/accounts/{tenant}/oauth2/auth [get]
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

	var client models.Client
	if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	if err := service.IsValidateUri(tenant.Id, client.Id, redirectUri); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid redirect_uri."})
		global.LOG.Error("get redirect uri err: " + err.Error())
		return
	}

	if responseType == "code" {
		code, err := getAccessCode(c, &client)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			global.LOG.Error("get access code err: " + err.Error())
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
		token, err := internal.GetAccessToken(c, &client)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			global.LOG.Error("get access token err: " + err.Error())
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

	c.String(http.StatusBadRequest, "Invalid response_type.")
	fmt.Println(clientId, scope, responseType, redirectUri, state, nonce)
}

// GetDeviceCode godoc
//
//	@Summary	device code
//	@Schemes
//	@Description	delete device groups
//	@Tags			oauth2
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client_id	query	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/oauth2/device/code [post]
func GetDeviceCode(c *gin.Context) {
	clientId := c.Query("client_id")
	scope := c.Query("scope")
	tenant := internal.GetTenant(c)
	var client models.Client
	if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
		global.LOG.Error("get client err: " + err.Error())
		return
	}

	deviceCode := models.DeviceCode{
		TenantId: tenant.Id,
		Code:     uuid.NewString(),
		UserCode: utils.GetDeviceUserCode(),
		Status:   "verifying",
	}

	if err := internal.TenantDB(c).Create(&deviceCode).Error; err != nil {
		c.String(http.StatusInternalServerError, "failed to create device code")
		global.LOG.Error("create device code err: " + err.Error())
		return
	}

	verificationUri := utils.GetHostWithScheme(c) + "/accounts/admin/" + c.Param("tenant") + "/devices/code"
	c.JSON(http.StatusOK, &gin.H{
		"deviceCode":              deviceCode.Code,
		"userCode":                deviceCode.UserCode,
		"clientId":                clientId,
		"expires_in":              2 * 60,
		"scope":                   scope,
		"verificationUri":         verificationUri,
		"verificationUriComplete": verificationUri + "/" + deviceCode.UserCode,
	})
}

// GetToken godoc
//
//	@Summary	oauth2 token
//	@Schemes
//	@Description	oauth2 token
//	@Tags			oauth2
//	@Param			tenant			path		string	true	"tenant"
//	@Param			client_id		query		string	true	"client_id"
//	@Param			client_secret	query		string	false	"client_secret"
//	@Param			code			query		string	false	"code"
//	@Param			grant_type		query		string	true	"grant_type"
//	@Param			redirect_uri	query		string	false	"redirect_uri"
//	@Param			state			query		string	false	"state"
//	@Param			nonce			query		string	false	"nonce"
//	@Success		200				{object}	dto.AccessTokenDto
//	@Router			/accounts/{tenant}/oauth2/token [get]
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
		var client models.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
			global.LOG.Error("get client err: " + err.Error())
			return
		}
		var secret models.ClientSecret
		if err := internal.TenantDB(c).First(&secret, "client_id = ? AND secret = ?", client.Id, clientSecret).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_secret."})
			global.LOG.Error("get client secret err: " + err.Error())
			return
		}

		var tokenCode models.TokenCode
		if err := internal.TenantDB(c).First(&tokenCode, "code = ?", code).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid code."})
			global.LOG.Error("get token code err: " + err.Error())
			return
		}
		service.ClearTokenCode(tokenCode.Code)
		accessToken := dto.AccessTokenDto{AccessToken: tokenCode.Token}
		c.JSON(http.StatusOK, accessToken)
		return
	} else if grantType == "client_credential" {
		tenant := internal.GetTenant(c)
		var client models.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
			global.LOG.Error("get client err: " + err.Error())
			return
		}
		var secret models.ClientSecret
		if err := internal.TenantDB(c).First(&secret, "client_id = ? AND secret = ?", client.Id, clientSecret).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_secret."})
			global.LOG.Error("get client secret err: " + err.Error())
			return
		}

		token, err := internal.GetClientAccessToken(c, &client)
		if err != nil {
			global.LOG.Error("get accessToken err: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "generate token err"})
			return
		}
		c.JSON(http.StatusOK, dto.AccessTokenDto{AccessToken: token})
		return
	} else if grantType == "urn:ietf:params:oauth:grant-type:device_code" {
		tenant := internal.GetTenant(c)
		var client models.Client
		if err := global.DB.First(&client, "tenant_id = ? AND id = ?", tenant.Id, clientId).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
			global.LOG.Error("get client err: " + err.Error())
			return
		}

		var deviceCode models.DeviceCode
		if err := global.DB.First(&deviceCode, "tenant_id = ? AND user_code = ?", tenant.Id, code).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid user code."})
			global.LOG.Error("get device code err: " + err.Error())
			return
		}
		if deviceCode.Status != "verified" {
			c.JSON(http.StatusForbidden, gin.H{"message": "device code is unauthorized."})
			return
		}

		token, err := internal.GetClientAccessToken(c, &client)
		if err != nil {
			global.LOG.Error("get accessToken err: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "generate token err"})
			return
		}

		service.ClearDeviceCode(deviceCode.UserCode)
		c.JSON(http.StatusOK, dto.AccessTokenDto{AccessToken: token})
		return
	} else if grantType == "device_credential" {
		id := c.Query("device_id")
		secret := c.Query("device_secret")
		var device models.Device
		if err := internal.TenantDB(c).Where("id = ?", id).First(&device).Error; err != nil {
			c.String(http.StatusUnauthorized, "invalidate device id")
			global.LOG.Error("get device id err: " + err.Error())
			return
		}

		var deviceSecret models.DeviceSecret
		if err := internal.TenantDB(c).Where("device_id = ? AND secret = ?", id, secret).First(&deviceSecret).Error; err != nil {
			c.String(http.StatusUnauthorized, "invalidate device secret")
			global.LOG.Error("get device secret err: " + err.Error())
			return
		}

		token, err := internal.GetDeviceToken(c, &device)
		if err != nil {
			global.LOG.Error("get accessToken err: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "generate token err"})
			return
		}

		c.JSON(http.StatusOK, dto.AccessTokenDto{AccessToken: token})
		return
	}

	c.String(http.StatusBadRequest, "Invalid grant_type.")
	fmt.Println(clientId, clientSecret, grantType, code, redirectUri, state, nonce)
}

// GetOpenidConfiguration godoc
//
//	@Summary	openid configuration
//	@Schemes
//	@Description	openid configuration
//	@Tags			oauth2
//	@Param			tenant			path		string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/.well-known/openid-configuration [get]
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
	c.JSON(http.StatusOK, conf)
}

// GetJwks godoc
//
//	@Summary	jwk
//	@Schemes
//	@Description	jwk
//	@Tags			oauth2
//	@Param			tenant			path		string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/.well-known/jwks.json [get]
func GetJwks(c *gin.Context) {
	tenant := internal.GetTenant(c)
	jwks, err := utils.LoadRsaPublicKeys(tenant.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get jwks err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, jwks)
}

func addOAuth2Routes(rg *gin.RouterGroup) {
	rg.GET("/oauth2/auth", middlewares.Authorized(true), GetAuthCode)
	rg.POST("/oauth2/device/code", GetDeviceCode)
	rg.GET("/oauth2/token", GetToken)
	rg.GET("/.well-known/openid-configuration", GetOpenidConfiguration)
	rg.GET("/.well-known/jwks.json", GetJwks)
}

package routes

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"accounts/utils"
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAccessToken(c *gin.Context, client *models.Client) (string, error) {
	user := GetUser(c)
	tenant := middlewares.GetTenant(c)
	scope := c.Query("scope")
	var clientUser models.ClientUser
	if err := data.DB.First(&clientUser, "tenant_id = ? AND client_id = ? AND user_id = ?", client.TenantId, client.Id, user.Id).Error; err != nil {
		clientUser.TenantId = client.TenantId
		clientUser.ClientId = client.Id
		clientUser.UserId = user.Id
		clientUser.Sub = uuid.NewString()
		if err := data.DB.Create(&clientUser).Error; err != nil {
			return "", err
		}
	}

	iss := fmt.Sprintf("%s/%s", utils.GetHostWithScheme(c), tenant.Name)
	now := time.Now()

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["sub"] = clientUser.Sub
	claims["aud"] = []string{client.CliId}
	claims["azp"] = client.CliId
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["name"] = user.Name()
	claims["scope"] = scope

	keys, err := utils.LoadRsaPrivateKeys(tenant.Name)
	if err != nil {
		return "", nil
	}
	var kid string
	var key *rsa.PrivateKey
	for kid, key = range keys {
		break
	}
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// clearTokenCode Deleted expired codes AND specific code.
func clearTokenCode(code string) {
	earliest := time.Now().Add(-2 * time.Minute)
	if err := data.DB.Delete(&models.TokenCode{}, "code = ? OR created_at < ?", code, earliest); err != nil {
		println(err)
	}
}

func getAccessCode(c *gin.Context, client *models.Client) (string, error) {
	token, err := getAccessToken(c, client)
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
	if err := data.DB.Create(&tokenCode).Error; err != nil {
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

	tenant := middlewares.GetTenant(c)
	var client models.Client
	if data.DB.First(&client, "tenant_id = ? AND client_id = ?", tenant.Id, clientId).Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
		return
	}
	var uri models.RedirectUri
	if data.DB.First(&uri, "tenant_id = ? AND client_id = ? AND redirect_uri = ?", tenant.Id, client.Id, redirectUri).Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid redirect_uri."})
		return
	}
	if responseType == "code" {
		code, err := getAccessCode(c, &client)
		if err != nil {
			c.Status(http.StatusInternalServerError)
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
	}

	if responseType == "token" {
		token, err := getAccessToken(c, &client)
		if err != nil {
			c.Status(http.StatusInternalServerError)
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

// GetToken godoc
//
//	@Summary	oauth2 token
//	@Schemes
//	@Description	oauth2 token
//	@Tags			oauth2
//	@Param			tenant			path		string	true	"tenant"
//	@Param			client_id		formData	string	true	"client_id"
//	@Param			client_secret	formData	string	true	"client_secret"
//	@Param			code			formData	string	true	"code"
//	@Param			scope			formData	string	true	"scope"
//	@Param			grant_type		formData	string	true	"grant_type"
//	@Param			redirect_uri	formData	string	true	"redirect_uri"
//	@Param			state			formData	string	false	"state"
//	@Param			nonce			formData	string	false	"nonce"
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
		tenant := middlewares.GetTenant(c)
		var client models.Client
		if data.DB.First(&client, "tenant_id = ? AND client_id = ?", tenant.Id, clientId).Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_id."})
			return
		}
		var secret models.ClientSecret
		if middlewares.TenantDB(c).First(&secret, "client_id = ? AND secret = ?", client.Id, clientSecret).Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid client_secret."})
			return
		}

		var tokenCode models.TokenCode
		if middlewares.TenantDB(c).First(&tokenCode, "code = ?", code).Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid code."})
			return
		}
		clearTokenCode(tokenCode.Code)
		accessToken := dto.AccessTokenDto{AccessToken: tokenCode.Token}
		c.JSON(http.StatusOK, accessToken)
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
//	@Success		200
//	@Router			/accounts/{tenant}/.well-known/openid-configuration [get]
func GetOpenidConfiguration(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
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
//	@Success		200
//	@Router			/accounts/.well-known/jwks.json [get]
func GetJwks(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	jwks, err := utils.LoadKeys(tenant.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, jwks)
}

func addOAuth2Routes(rg *gin.RouterGroup) {
	rg.GET("/oauth2/auth", middlewares.Authorized(true), GetAuthCode)
	rg.GET("/oxauth2/token", middlewares.Authorized(false), GetToken)
	rg.GET("/.well-known/openid-configuration", GetOpenidConfiguration)
	rg.GET("/.well-known/jwks.json", GetJwks)
}

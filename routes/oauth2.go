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

func GetHostWithScheme(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if s := c.Request.Header.Get("X-Forwarded-Proto"); s != "" {
		scheme = s
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

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

	iss := fmt.Sprintf("%s/%s", GetHostWithScheme(c), tenant.Name)
	now := time.Now()

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["sub"] = clientUser.Sub
	claims["aud"] = []string{client.ClientId}
	claims["azp"] = client.ClientId
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
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

func getAccessTokenFromCode(c *gin.Context, client *models.Client) {

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

func addOAuth2Routes(rg *gin.RouterGroup) {
	// If logged in, 302 to redirect uri.
	// If not, return login form.
	rg.GET("/oauth2/auth", middlewares.Authorized, func(c *gin.Context) {
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
	})

	rg.GET("/oauth2/token", middlewares.Authorized, func(c *gin.Context) {
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
			c.JSON(http.StatusOK, gin.H{"access_token": tokenCode.Token})
			return
		}

		c.String(http.StatusBadRequest, "Invalid grant_type.")
		fmt.Println(clientId, clientSecret, grantType, code, redirectUri, state, nonce)
	})

	rg.GET("/.well-known/openid-configuration", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		prefix := GetHostWithScheme(c)
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
	})

	rg.GET("/.well-known/jwks.json", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		jwks, err := utils.LoadKeys(tenant.Name)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, jwks)
	})
}

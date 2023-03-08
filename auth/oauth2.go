package auth

import (
	"accounts/models"
	"accounts/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"net/url"
)

type ProviderOAuth2 struct {
	Id         uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint            `json:"providerId"`
	Provider   models.Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	TokenEndpoint     string `json:"tokenEndpoint"`
	UserinfoEndpoint  string `json:"userinfoEndpoint"`
	Scope             string `json:"scope"`
	ResponseType      string `json:"responseType"`

	TenantId uint `gorm:"primaryKey"`
}

func (ProviderOAuth2) TableName() string {
	return "provider_oauth2"
}

func (p ProviderOAuth2) Auth(redirectUri string) string {
	query := url.Values{}
	query.Set("client_id", p.ClientId)
	query.Set("scope", p.Scope)
	query.Set("response_type", p.ResponseType)
	query.Set("redirect_uri", redirectUri)
	location := fmt.Sprintf("%s?%s", p.AuthorizeEndpoint, query.Encode())
	return location
}

func (p ProviderOAuth2) Login(c *gin.Context) (*UserInfo, error) {
	tenantName := c.Param("tenant")
	providerName := c.Param("provider")
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenantName, providerName)
	query := url.Values{}
	query.Set("client_id", p.ClientId)
	query.Set("client_secret", p.ClientSecret)
	query.Set("scope", p.Scope)
	query.Set("code", code)
	query.Set("redirect_uri", redirectUri)
	query.Set("grant_type", "authorization_code")
	resp, err := http.PostForm(p.TokenEndpoint, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	tokenString := result["access_token"]
	if tokenString == nil {
		return nil, err
	}

	token, _ := jwt.Parse(tokenString.(string), nil)
	if token == nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)

	return &UserInfo{
		Sub:         GetString(claims["sub"]),
		DisplayName: GetString(claims["name"]),
		FirstName:   GetString(claims["given_name"]),
		LastName:    GetString(claims["family_name"]),
		Email:       GetString(claims["email"]),
		Phone:       GetString(claims["phone_number"]),
	}, nil
}

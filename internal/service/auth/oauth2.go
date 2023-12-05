package auth

import (
	"alfred/internal/model"
	"alfred/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/url"
)

type ProviderOAuth2 struct {
	Config model.ProviderOAuth2
}

func (ProviderOAuth2) TableName() string {
	return "provider_oauth2"
}

func (p ProviderOAuth2) Auth(redirectUri string, state string, _ uint) (string, error) {
	query := url.Values{}
	query.Set("client_id", p.Config.ClientId)
	query.Set("scope", p.Config.Scope)
	query.Set("response_type", p.Config.ResponseType)
	query.Set("redirect_uri", redirectUri)
	query.Set("state", state)
	location := fmt.Sprintf("%s?%s", p.Config.AuthorizeEndpoint, query.Encode())
	return location, nil
}

func (p ProviderOAuth2) Login(c *gin.Context) (*model.UserInfo, error) {
	tenantName := c.Param("tenant")
	providerName := c.Param("provider")
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenantName, providerName)
	query := url.Values{}
	query.Set("client_id", p.Config.ClientId)
	query.Set("client_secret", p.Config.ClientSecret)
	query.Set("scope", p.Config.Scope)
	query.Set("code", code)
	query.Set("redirect_uri", redirectUri)
	query.Set("grant_type", "authorization_code")
	resp, err := http.PostForm(p.Config.TokenEndpoint, query)
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

	return &model.UserInfo{
		Sub:         utils.GetString(claims["sub"]),
		DisplayName: utils.GetString(claims["name"]),
		FirstName:   utils.GetString(claims["given_name"]),
		LastName:    utils.GetString(claims["family_name"]),
		Email:       utils.GetString(claims["email"]),
		Phone:       utils.GetString(claims["phone_number"]),
	}, nil
}

func (p ProviderOAuth2) LoginConfig() *gin.H {
	return &gin.H{
		"type":         p.Config.Provider.Type,
		"providerId":   p.Config.ProviderId,
		"clientId":     p.Config.ClientId,
		"responseType": p.Config.ResponseType,
		"scope":        p.Config.Scope,
	}
}

func (p ProviderOAuth2) ProviderConfig() *gin.H {
	return &gin.H{
		"type":         p.Config.Provider.Type,
		"providerId":   p.Config.ProviderId,
		"clientId":     p.Config.ClientId,
		"clientSecret": p.Config.ClientSecret,
		"responseType": p.Config.ResponseType,
		"scope":        p.Config.Scope,
	}
}

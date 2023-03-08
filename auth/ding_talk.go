package auth

import (
	"accounts/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type ProviderDingTalk struct {
	Id         uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint            `json:"providerId"`
	Provider   models.Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	AgentId   string
	AppKey    string
	AppSecret string

	TenantId uint `gorm:"primaryKey"`
}

func (p ProviderDingTalk) Auth(redirectUri string) string {
	query := url.Values{}
	query.Set("client_id", p.AppKey)
	query.Set("scope", "openid corpid")
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectUri)
	query.Set("state", uuid.NewString())
	query.Set("prompt", "consent")
	location := fmt.Sprintf("%s?%s", "https://login.dingtalk.com/oauth2/auth", query.Encode())
	return location
}

type dingTalkTokenRequest struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrantType    string `json:"grantType"`
}

func (p ProviderDingTalk) Login(c *gin.Context) (*UserInfo, error) {
	code := c.Query("authCode")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	request := dingTalkTokenRequest{
		ClientId:     p.AppKey,
		ClientSecret: p.AppSecret,
		Code:         code,
		GrantType:    "authorization_code",
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://api.dingtalk.com//v1.0/oauth2/userAccessToken", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	accessToken := GetString(result["accessToken"])
	if accessToken == "" {
		return nil, err
	}

	// Get Profile.
	url := "https://api.dingtalk.com/v1.0/contact/users/me"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var profile map[string]any
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return nil, err
	}
	if GetString(profile["openId"]) == "" {
		return nil, errors.New("get userinfo failed")
	}

	return &UserInfo{
		Sub:         GetString(profile["openId"]),
		DisplayName: GetString(profile["nick"]),
		Email:       GetString(profile["email"]),
		Phone:       GetString(profile["mobile"]),
		Picture:     GetString(profile["avatarUrl"]),
	}, nil
}

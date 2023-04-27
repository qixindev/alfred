package auth

import (
	"accounts/models"
	"accounts/utils"
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
	Config models.ProviderDingTalk
}

func (p ProviderDingTalk) Auth(redirectUri string) (string, error) {
	query := url.Values{}
	query.Set("client_id", p.Config.AppKey)
	query.Set("scope", "openid corpid")
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectUri)
	query.Set("state", uuid.NewString())
	query.Set("prompt", "consent")
	location := fmt.Sprintf("%s?%s", "https://login.dingtalk.com/oauth2/auth", query.Encode())
	return location, nil
}

type dingTalkTokenRequest struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrantType    string `json:"grantType"`
}

func (p ProviderDingTalk) Login(c *gin.Context) (*models.UserInfo, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	request := dingTalkTokenRequest{
		ClientId:     p.Config.AppKey,
		ClientSecret: p.Config.AppSecret,
		Code:         code,
		GrantType:    "authorization_code",
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	
	resp, err := http.Post("https://api.dingtalk.com/v1.0/oauth2/userAccessToken", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result map[string]any
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	accessToken := utils.GetString(result["accessToken"])
	if accessToken == "" {
		return nil, errors.New("failed to get ding token")
	}

	// Get Profile.
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.dingtalk.com/v1.0/contact/users/me", nil)
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
	if err = json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return nil, err
	}

	if utils.GetString(profile["openId"]) == "" {
		return nil, errors.New("get userinfo failed")
	}

	return &models.UserInfo{
		Sub:         utils.GetString(profile["openId"]),
		DisplayName: utils.GetString(profile["nick"]),
		Email:       utils.GetString(profile["email"]),
		Phone:       utils.GetString(profile["mobile"]),
		Picture:     utils.GetString(profile["avatarUrl"]),
	}, nil
}

func (p ProviderDingTalk) LoginConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
		"appKey":     p.Config.AppKey,
		"type":       p.Config.Provider.Type,
	}
}

package auth

import (
	"accounts/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type ProviderWeCom struct {
	Id         uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint            `json:"providerId"`
	Provider   models.Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	CorpId    string `json:"corpId"`
	AgentId   string `json:"agentId"`
	AppSecret string `json:"appSecret"`

	TenantId uint `gorm:"primaryKey"`
}

func (p ProviderWeCom) Auth(redirectUri string) string {
	query := url.Values{}
	query.Set("appid", p.CorpId)
	query.Set("scope", "snsapi_base")
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectUri)
	query.Set("agentid", p.AgentId)
	query.Set("state", uuid.NewString())
	location := fmt.Sprintf("%s?%s#wechat_redirect", "https://open.weixin.qq.com/connect/oauth2/authorize", query.Encode())
	return location
}

func (p ProviderWeCom) Login(c *gin.Context) (*UserInfo, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	query := url.Values{}
	query.Set("corpid", p.CorpId)
	query.Set("corpsecret", p.AppSecret)

	resp, err := http.Get(fmt.Sprintf("%s?%s", "https://qyapi.weixin.qq.com/cgi-bin/gettoken", query.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	accessToken := GetString(result["access_token"])
	if accessToken == "" {
		return nil, err
	}

	// Basic info.
	basicInfoQuery := url.Values{}
	basicInfoQuery.Set("access_token", accessToken)
	basicInfoQuery.Set("code", code)
	basicInfoUrl := "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"
	resp, err = http.Get(fmt.Sprintf("%s?%s", basicInfoUrl, basicInfoQuery.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	userId := GetString(result["userid"])
	if userId == "" {
		return nil, err
	}
	var userInfo UserInfo
	userInfo.Sub = userId
	userInfo.DisplayName = GetString(result["userid"])

	detailInfoUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	detailInfoQuery := url.Values{}
	detailInfoQuery.Set("access_token", accessToken)
	detailInfoQuery.Set("userid", userId)
	resp, err = http.Get(fmt.Sprintf("%s?%s", detailInfoUrl, detailInfoQuery.Encode()))
	if err != nil {
		return &userInfo, nil
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return &userInfo, nil
	}
	userInfo.DisplayName = GetString(result["name"])

	return &userInfo, nil
}

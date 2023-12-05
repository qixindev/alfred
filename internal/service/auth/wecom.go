package auth

import (
	"alfred/internal/model"
	"alfred/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type ProviderWeCom struct {
	Config model.ProviderWeCom
}

func (p ProviderWeCom) Auth(redirectUri string, state string, _ uint) (string, error) {
	query := url.Values{}
	query.Set("appid", p.Config.CorpId)
	query.Set("scope", "snsapi_base")
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectUri)
	query.Set("agentid", p.Config.AgentId)
	query.Set("state", state)
	location := fmt.Sprintf("%s?%s#wechat_redirect", "https://open.weixin.qq.com/connect/oauth2/authorize", query.Encode())
	return location, nil
}

func (p ProviderWeCom) Login(c *gin.Context) (*model.UserInfo, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	query := url.Values{}
	query.Set("corpid", p.Config.CorpId)
	query.Set("corpsecret", p.Config.AppSecret)
	resp, err := http.Get(fmt.Sprintf("%s?%s", "https://qyapi.weixin.qq.com/cgi-bin/gettoken", query.Encode()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result map[string]any
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	accessToken := utils.GetString(result["access_token"])
	if accessToken == "" {
		return nil, errors.New("failed to get ding token")
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
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	userId := utils.GetString(result["userid"])
	if userId == "" {
		return nil, err
	}

	var userInfo model.UserInfo
	userInfo.Sub = userId
	userInfo.DisplayName = utils.GetString(result["userid"])
	detailInfoUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	detailInfoQuery := url.Values{}
	detailInfoQuery.Set("access_token", accessToken)
	detailInfoQuery.Set("userid", userId)
	resp, err = http.Get(fmt.Sprintf("%s?%s", detailInfoUrl, detailInfoQuery.Encode()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	userInfo.DisplayName = utils.GetString(result["name"])

	return &userInfo, nil
}

func (p ProviderWeCom) LoginConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
		"corpId":     p.Config.CorpId,
		"agentId":    p.Config.AgentId,
		"type":       p.Config.Provider.Type,
	}
}

func (p ProviderWeCom) ProviderConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
		"corpId":     p.Config.CorpId,
		"agentId":    p.Config.AgentId,
		"appSecret":  p.Config.AppSecret,
		"type":       p.Config.Provider.Type,
	}
}

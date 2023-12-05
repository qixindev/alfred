package auth

import (
	"alfred/internal/model"
	"alfred/pkg/client/msg/api"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
)

type ProviderWechat struct {
	Config model.ProviderWechat
}

func (p ProviderWechat) Auth(redirectUri string, state string, _ uint) (string, error) {
	query := url.Values{}
	query.Set("appid", p.Config.AppId)
	query.Set("redirect_uri", redirectUri)
	query.Set("response_type", "code")
	query.Set("scope", "snsapi_login")
	query.Set("state", state)
	location := fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?%s#wechat_redirect", query.Encode())
	return location, nil
}

type WechatTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
	ErrorCode    int    `json:"errcode"`
	ErrorMsg     string `json:"errmsg"`
}
type WechatUserInfo struct {
	Openid     string        `json:"openid"`
	Nickname   string        `json:"nickname"`
	Sex        int           `json:"sex"`
	Language   string        `json:"language"`
	City       string        `json:"city"`
	Province   string        `json:"province"`
	Country    string        `json:"country"`
	Headimgurl string        `json:"headimgurl"`
	Privilege  []interface{} `json:"privilege"`
	Unionid    string        `json:"unionid"`
	ErrorCode  int           `json:"errcode"`
	ErrorMsg   string        `json:"errmsg"`
}

func (p ProviderWechat) Login(c *gin.Context) (*model.UserInfo, error) {
	code := c.Query("code")
	if code == "" {
		return nil, errors.New("no auth code")
	}
	query := url.Values{}
	query.Set("appid", p.Config.AppId)
	query.Set("secret", p.Config.AppSecret)
	query.Set("code", code)
	query.Set("grant_type", "authorization_code")
	tokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?" + query.Encode()
	var t WechatTokenResp
	if err := api.GetClient(tokenUrl, &t); err != nil {
		return nil, errors.WithMessage(err, "failed to get wechat access token")
	}
	if t.ErrorCode != 0 || t.ErrorMsg != "" {
		return nil, errors.New(fmt.Sprintf("get wechat access token err: %d:%s", t.ErrorCode, t.ErrorMsg))
	}
	userInfoQuery := url.Values{}
	userInfoQuery.Set("access_token", t.AccessToken)
	userInfoQuery.Set("openid", t.Openid)
	userInfoUrl := "https://api.weixin.qq.com/sns/userinfo?" + userInfoQuery.Encode()
	var wechatUserInfo WechatUserInfo
	if err := api.GetClient(userInfoUrl, &wechatUserInfo); err != nil {
		return nil, errors.WithMessage(err, "failed to get wechat user info")
	}

	if wechatUserInfo.ErrorCode != 0 || wechatUserInfo.ErrorMsg != "" {
		return nil, errors.New(fmt.Sprintf("get wechat access token err: %d:%s", t.ErrorCode, t.ErrorMsg))
	}
	userInfo := model.UserInfo{
		Sub:         t.Openid,
		DisplayName: wechatUserInfo.Nickname,
		Picture:     wechatUserInfo.Headimgurl,
	}

	global.LOG.Info("wechat user info: " + utils.StructToString(userInfo))
	return &userInfo, nil
}

func (p ProviderWechat) LoginConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
		"appId":      p.Config.AppId,
		"type":       p.Config.Provider.Type,
	}
}

func (p ProviderWechat) ProviderConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
		"appId":      p.Config.AppId,
		"appSecret":  p.Config.AppSecret,
		"type":       p.Config.Provider.Type,
	}
}

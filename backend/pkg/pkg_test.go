package pkg

import (
	"alfred/backend/pkg/client/msg/api"
	"alfred/backend/pkg/utils"
	"fmt"
	"net/url"
	"testing"
)

type WechatTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
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
}

func TestTemp(t *testing.T) {
	query := url.Values{}
	query.Set("appid", "wx7f8b2509c7f46c05")
	query.Set("secret", "cd2f96d02820bc013829d87e51acd2c0")
	query.Set("code", "02152s100HyA9R1vIb400vRQ20052s1R")
	query.Set("grant_type", "authorization_code")
	tokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?" + query.Encode()
	var tt WechatTokenResp
	if err := api.GetClient(tokenUrl, &tt); err != nil {
		fmt.Println("token err: ", err)
		return
	}

	fmt.Println(utils.StructToString(tt))
	userInfoQuery := url.Values{}
	userInfoQuery.Set("access_token", tt.AccessToken)
	userInfoQuery.Set("openid", tt.Openid)
	userInfoUrl := "https://api.weixin.qq.com/sns/userinfo?" + userInfoQuery.Encode()
	var wechatUserInfo WechatUserInfo
	if err := api.GetClient(userInfoUrl, &wechatUserInfo); err != nil {
		fmt.Println("userinfo err: ", err)
		return
	}

	fmt.Println(utils.StructToString(wechatUserInfo))
}

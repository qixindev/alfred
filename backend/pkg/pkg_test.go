package pkg

import (
	"alfred/backend/pkg/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Get("http://localhost/accounts/default/providers/sms/login?phone=%2B8613365807972")
	if err != nil {
		t.FailNow()
	}
	defer utils.DeferErr(resp.Body.Close)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.FailNow()
	}
	result := struct {
		Location string `json:"location"`
		State    string `json:"state"`
	}{}
	if err = json.Unmarshal(body, &result); err != nil {
		t.FailNow()
	}

	println(utils.StructToString(result))
	url := fmt.Sprintf("http://localhost/accounts/login/providers/callback?state=%s&code=%s", result.State, result.Location)
	resp, err = client.Get(url)
	if err != nil {
		t.FailNow()
	}
	defer utils.DeferErr(resp.Body.Close)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.FailNow()
	}
	println(string(body))
}

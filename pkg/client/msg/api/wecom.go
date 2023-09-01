package api

import (
	"errors"
	"fmt"
	"strconv"
)

type MarkdownContext struct {
	Content string `json:"content"`
}
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PicUrl      string `json:"picurl"`
	Appid       string `json:"appid"`
	PagePath    string `json:"pagepath"`
}
type NewsMsg struct {
	Articles []Article `json:"articles"`
}
type MsgStruct struct {
	ToUser               string          `json:"touser"`
	ToParty              string          `json:"toparty"`
	ToTag                string          `json:"totag"`
	MsgType              string          `json:"msgtype"`
	AgentId  int64           `json:"agentid"`
	Markdown MarkdownContext `json:"markdown"`
	News     NewsMsg         `json:"news"`
	Safe     int             `json:"safe"`
	EnableIdTrans        int             `json:"enable_id_trans"`
	EnableDuplicateCheck int             `json:"enable_duplicate_check"`
}

type AccessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func ForceGetWecomToken(corpId string, corpSecret string) (string, error) {
	var url = fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, corpSecret)
	var res AccessToken
	if err := PostClient(url, struct{}{}, &res); err != nil {
		return "", err
	}

	if res.ErrCode != 0 && res.ErrMsg != "ok" || res.AccessToken == "" {
		return "", errors.New("Get token err: [" + strconv.Itoa(res.ErrCode) + "]" + res.ErrMsg)
	}

	return res.AccessToken, nil
}

func SendWecomMsg(token string, markdown MsgStruct) error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	var resp struct {
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
		InvalidParty string `json:"invalidparty"`
		InvalidTag   string `json:"invalidtag"`
		MsgId        string `json:"msgid"`
	}
	if err := PostClient(url, markdown, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 && resp.ErrMsg != "ok" {
		return errors.New("Send msg err: [" + strconv.Itoa(resp.ErrCode) + "] " + resp.ErrMsg)
	}

	return nil
}

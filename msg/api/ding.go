package api

import (
	"errors"
	"fmt"
)

type DingMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type btnJsonList struct {
	Title     string `json:"title"`
	ActionUrl string `json:"action_url"`
}
type DingActionCard struct {
	Title          string        `json:"title"`
	Markdown       string        `json:"markdown"`
	BtnOrientation string        `json:"btn_orientation"`
	BtnJsonList    []btnJsonList `json:"btn_json_list"`
	SingleTitle    string        `json:"single_title"`
	SingleUrl      string        `json:"single_url"`
}

type DingMessage struct {
	MsgType    string         `json:"msgtype"`
	ActionCard DingActionCard `json:"action_card"`
	Markdown   DingMarkdown   `json:"markdown"`
}

type DingNotify struct {
	AgentId    int64       `json:"agent_id"`
	UseridList string      `json:"userid_list"`
	ToAllUser  bool        `json:"to_all_user"`
	Msg        DingMessage `json:"msg"`
}

func GetDingAccessToken(conf *Ding) (string, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", conf.AppKey, conf.AppSecret)

	var resp struct {
		ErrCode float64 `json:"errcode"`
		ErrMsg  string  `json:"errmsg"`
		Token   string  `json:"access_token"`
	}
	if err := GetClient(url, &resp); err != nil {
		return "", err
	}

	if resp.ErrCode != 0 || resp.Token == "" {
		return "", errors.New(fmt.Sprintf("token获取错误: [%f] %s", resp.ErrCode, resp.ErrMsg))
	}

	return resp.Token, nil
}

func getSendMsgRes(taskId int64, conf *Ding) error {
	token, err := GetDingAccessToken(conf)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/message/corpconversation/getsendresult?access_token=%s", token)
	body := struct {
		AgentId int64 `json:"agent_id"`
		TaskId  int64 `json:"task_id"`
	}{conf.AgentId, taskId}
	var resp struct {
		ErrCode float64 `json:"errcode"`
		ErrMsg  string  `json:"errmsg"`
	}
	if err = PostClient(url, body, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 || resp.ErrMsg != "ok" {
		return errors.New(fmt.Sprintf("发送钉钉消息失败: [%f] %s", resp.ErrCode, resp.ErrMsg))
	}

	return nil
}

func GetUseridByUnionId(token, unionId string) (userid string, err error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=%s", token)
	p := struct {
		UnionId string `json:"unionid"`
	}{unionId}

	var resp struct {
		ErrCode float64 `json:"errcode"`
		ErrMsg  string  `json:"errmsg"`
		Result  struct {
			UserId string `json:"userid"`
		}
	}
	if err = PostClient(url, p, &resp); err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", errors.New(fmt.Sprintf("userid获取错误: [%f] %s", resp.ErrCode, resp.ErrMsg))
	}

	return resp.Result.UserId, nil
}

// SendDingMsg 发送消息到钉钉
func SendDingMsg(token string, dingMsg DingNotify, conf *Ding) error {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s", token)
	var resp struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		TaskId    int64  `json:"task_id"`
		RequestId string `json:"request_id"`
	}

	if err := PostClient(url, dingMsg, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 || resp.TaskId == 0 {
		return errors.New(fmt.Sprintf("发送钉钉消息失败: [%d] %s", resp.ErrCode, resp.ErrMsg))
	}

	if err := getSendMsgRes(resp.TaskId, conf); err != nil {
		return err
	}

	return nil
}

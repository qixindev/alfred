package notify

import (
	"accounts/config/env"
	"accounts/msg/api"
	"accounts/utils"
	"net/url"
)

func getMarkdownString(info *SendInfo) string {
	text := info.Msg
	if info.Title != "" {
		text = "<font color=" + info.TitleColor + ">" + info.Title + "</font>" + text
	}

	if info.PngLink != "" {
		text = "![logo](" + info.PngLink + ")   \n" + text
	}

	return text
}

// markdown消息
func getMarkDownMsg(info *SendInfo) api.DingMessage {
	markdownMsg := api.DingMessage{
		MsgType: "markdown",
		Markdown: api.DingMarkdown{
			Title: "通知",
			Text:  getMarkdownString(info),
		},
	}
	return markdownMsg
}

// 卡片消息
func getActionMsg(info *SendInfo) api.DingMessage {
	actionMsg := api.DingMessage{
		MsgType: "action_card",
		ActionCard: api.DingActionCard{
			Title:       info.Title,
			Markdown:    getMarkdownString(info),
			SingleTitle: "查看详情",
			SingleUrl:   "dingtalk://dingtalkclient/page/link?url=" + url.QueryEscape(info.Link) + "&pc_slide=false",
		},
	}
	return actionMsg
}

func GetDingMsg(info *SendInfo, conf *api.Third) api.DingNotify {
	dingMsg := api.DingNotify{
		AgentId:    conf.GetDingAgentId(),
		UseridList: utils.MergeString(info.Users, ","),
		ToAllUser:  false,
	}

	switch info.Type {
	case env.MsgMarkdown:
		dingMsg.Msg = getMarkDownMsg(info)
	default: // case config.MsgPicture:
		dingMsg.Msg = getActionMsg(info)
	}

	return dingMsg
}

func SendMsgToDingTalk(info *SendInfo, conf *api.Third) error {
	token, err := api.GetDingAccessToken(conf)
	if err != nil {
		return err
	}

	dingMsg := GetDingMsg(info, conf)
	return api.SendDingMsg(token, dingMsg, conf)
}

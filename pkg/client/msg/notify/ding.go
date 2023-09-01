package notify

import (
	"accounts/internal/model"
	"accounts/pkg/client/msg/api"
	"accounts/pkg/config/env"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
)

func getMarkdownString(info *model.SendInfo) string {
	text := info.Msg
	if info.Title != "" {
		text = "<font color=" + info.TitleColor + ">" + info.Title + "</font>   \n" + text
	}

	if info.PngLink != "" {
		text = "![logo](" + info.PngLink + ")   \n" + text
	}

	return text
}

// markdown消息
func getMarkDownMsg(info *model.SendInfo) api.DingMessage {
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
func getActionMsg(info *model.SendInfo) api.DingMessage {
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

func GetDingMsg(info *model.SendInfo, conf *api.Ding) api.DingNotify {
	dingMsg := api.DingNotify{
		AgentId:    conf.AgentId,
		UseridList: utils.MergeString(info.Users, ","),
		ToAllUser:  false,
	}

	switch info.MsgType {
	case env.MsgMarkdown:
		dingMsg.Msg = getMarkDownMsg(info)
	default: // case config.MsgPicture:
		dingMsg.Msg = getActionMsg(info)
	}

	return dingMsg
}

func convertToUserId(unionIds []string, token string) ([]string, error) {
	res := make([]string, 0, len(unionIds))
	for _, unionId := range unionIds {
		userId, err := api.GetUseridByUnionId(token, unionId)
		if err != nil {
			return nil, errors.Wrap(err, "get ding talk user id err")
		}
		res = append(res, userId)
	}
	return res, nil
}

func SendMsgToDingTalk(info *model.SendInfo, providerConf gin.H) error {
	conf, err := api.GetDingTalkConfig(providerConf)
	if err != nil {
		return err
	}
	token, err := api.GetDingAccessToken(conf)
	if err != nil {
		return err
	}

	info.Users, err = convertToUserId(info.Users, token)
	if err != nil {
		return err
	}
	dingMsg := GetDingMsg(info, conf)
	return api.SendDingMsg(token, dingMsg, conf)
}

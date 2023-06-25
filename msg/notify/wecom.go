package notify

import (
	"accounts/config/env"
	"accounts/msg/api"
	"accounts/utils"
)

func getMarkdownText(info *SendInfo) string {
	markdownText := "<font color=\"" + info.TitleColor + "\">" + info.Title + "</font>   \n" + info.Msg
	return markdownText
}

// 企业微信文章类消息
func getNews(title, text, link, logo string) api.NewsMsg {
	news := api.NewsMsg{
		Articles: []api.Article{{
			Title:       title,
			Description: text,
			Url:         link,
			PicUrl:      logo,
		}},
	}
	return news
}

func getMsg(info *SendInfo, conf *api.Third) api.MsgStruct {
	toUser := utils.MergeString(info.Users, "|")
	msg := api.MsgStruct{
		ToUser:               toUser,
		ToParty:              "@all",
		ToTag:                "@all",
		AgentId:              conf.GetWecomAgentId(),
		Safe:                 0,
		EnableIdTrans:        0,
		EnableDuplicateCheck: 0,
	}

	switch info.Type {
	case env.MsgMarkdown:
		msg.MsgType = "markdown"
		text := getMarkdownText(info)
		msg.Markdown = api.MarkdownContext{Content: text}
	default: // case MsgPicture:
		msg.MsgType = "news"
		msg.News = getNews(info.Title, info.Msg, info.Link, info.PngLink)
	}

	return msg
}

func SendMsgToWecom(info *SendInfo, conf *api.Third) error {
	corpId := conf.GetWecomCorpId()
	corpSecret := conf.GetWecomSecret()
	token, err := api.ForceGetWecomToken(corpId, corpSecret)
	if err != nil {
		return err
	}

	msgStruct := getMsg(info, conf)

	return api.SendWecomMsg(token, msgStruct)
}
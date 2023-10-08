package notify

import (
	"alfred/internal/model"
	"alfred/pkg/client/msg/api"
	"alfred/pkg/config/env"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
)

func getMarkdownText(info *model.SendInfo) string {
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

func getMsg(info *model.SendInfo, conf *api.Wecom) api.MsgStruct {
	toUser := utils.MergeString(info.Users, "|")
	msg := api.MsgStruct{
		ToUser:               toUser,
		ToParty:              "@all",
		ToTag:                "@all",
		AgentId:              conf.AgentId,
		Safe:                 0,
		EnableIdTrans:        0,
		EnableDuplicateCheck: 0,
	}

	switch info.MsgType {
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

func SendMsgToWecom(info *model.SendInfo, providerConfig gin.H) error {
	conf, err := api.GetWecomConfig(providerConfig)
	if err != nil {
		return err
	}

	token, err := api.ForceGetWecomToken(conf.CorpId, conf.Secret)
	if err != nil {
		return err
	}

	msgStruct := getMsg(info, conf)

	return api.SendWecomMsg(token, msgStruct)
}

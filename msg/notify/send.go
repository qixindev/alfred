package notify

import (
	"accounts/config/env"
	"errors"
	"github.com/gin-gonic/gin"
)

type SendInfo struct {
	Msg   string   // 要发送的消息
	Link  string   // 点击跳转链接
	Users []string // 发送给谁

	Platform   string // 发送到哪个平台
	MsgType    string // 消息类型：图文，markdown，文字
	Title      string // 标题
	TitleColor string // 标题颜色
	PngLink    string // 消息图片链接
}

func SendMsgToUsers(info *SendInfo, conf gin.H) (err error) {
	if len(info.Users) == 0 {
		return errors.New("users is null")
	}

	switch info.Platform {
	case env.PlatformDingTalk:
		err = SendMsgToDingTalk(info, conf)
	case env.PlatformWecom:
		err = SendMsgToWecom(info, conf)
	default:
		err = errors.New("platform not support: " + info.Platform)
	}
	return err
}

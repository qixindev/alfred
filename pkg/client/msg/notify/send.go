package notify

import (
	"accounts/internal/model"
	"accounts/pkg/config/env"
	"errors"
	"github.com/gin-gonic/gin"
)

func SendMsgToUsers(info *model.SendInfo, conf gin.H) (err error) {
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

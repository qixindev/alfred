package server

import (
	"accounts/global"
	"accounts/models"
	"accounts/msg/notify"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
)

// SendMsg godoc
//
//	@Summary	send message
//	@Schemes
//	@Description	send message
//	@Tags			msg
//	@Param			tenant		path	string	true	"tenant name"
//	@Param			providerId	path	string	true	"provider id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{providerId} [get]
func SendMsg(c *gin.Context) {
	var in notify.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	var provider models.Provider
	if err := global.DB.Where("").First(&provider).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to get provider info")
		global.LOG.Error("get send provider info err: " + err.Error())
		return
	}

	if err := notify.SendMsgToUsers(&in, nil); err != nil {
		internal.ErrorSqlResponse(c, "failed to send msg")
		global.LOG.Error("send msg err: " + err.Error())
		return
	}

	internal.Success(c)
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:providerId", SendMsg)
}

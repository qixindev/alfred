package server

import (
	"accounts/global"
	"accounts/msg/notify"
	"accounts/server/auth"
	"accounts/server/internal"
	"accounts/utils"
	"github.com/gin-gonic/gin"
)

// SendMsg godoc
//
//	@Summary	send message
//	@Schemes
//	@Description	send message
//	@Tags			msg
//	@Param			tenant		path	string			true	"tenant name"
//	@Param			providerId	path	string			true	"provider name"
//	@Param			by			body	notify.SendInfo	true	"msg body"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{provider} [get]
func SendMsg(c *gin.Context) {
	var in notify.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	tenant := internal.GetTenant(c)
	authProvider, err := auth.GetAuthProvider(tenant.Id, in.Platform)
	if err != nil {
		return
	}

	var providerUser []string
	providerConfig := *authProvider.ProviderConfig()
	if err = global.DB.Table("provider_users pu").Select("pu.name").
		Joins("LEFT JOIN client_users as cu ON cu.user_id = pu.user_id").
		Where("tenant_id = ? AND provider_id = ? AND cu.sub in ?", tenant.Id, providerConfig["providerId"], in.Users).
		Find(&providerUser).Error; err != nil {
		return
	}

	in.Users = providerUser
	global.LOG.Debug("msg send info: " + utils.StructToString(in))
	global.LOG.Debug("msg conf: " + utils.StructToString(providerConfig))
	if err = notify.SendMsgToUsers(&in, providerConfig); err != nil {
		internal.ErrorSqlResponse(c, "failed to send msg")
		global.LOG.Error("send msg err: " + err.Error())
		return
	}

	internal.Success(c)
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:providerId", SendMsg)
}

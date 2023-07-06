package server

import (
	"accounts/global"
	"accounts/msg/notify"
	"accounts/server/auth"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendMsg godoc
//
//	@Summary	send message
//	@Schemes
//	@Description	send message
//	@Tags			msg
//	@Param			tenant		path	string			true	"tenant name"
//	@Param			provider	path	string			true	"provider name"
//	@Param			by			body	notify.SendInfo	true	"msg body"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{provider} [post]
func SendMsg(c *gin.Context) {
	var in notify.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	tenant := internal.GetTenant(c)
	authProvider, err := auth.GetAuthProvider(tenant.Id, c.Param("provider"))
	if err != nil {
		global.LOG.Error("get provider err: " + err.Error())
		internal.ErrorSqlResponse(c, "no such provider")
		return
	}

	var providerUser []string
	providerConfig := *authProvider.ProviderConfig()
	if err = global.DB.Table("provider_users pu").Select("pu.name").
		Joins("LEFT JOIN client_users as cu ON cu.user_id = pu.user_id").
		Where("pu.tenant_id = ? AND pu.provider_id = ? AND cu.sub in ?", tenant.Id, providerConfig["providerId"], in.Users).
		Find(&providerUser).Error; err != nil {
		global.LOG.Error("get provider user err")
		internal.ErrorSqlResponse(c, "failed to get provider user")
		return
	}

	in.Users = providerUser
	in.Platform = providerConfig["type"].(string)
	if err = notify.SendMsgToUsers(&in, providerConfig); err != nil {
		global.LOG.Error("send msg err: " + err.Error())
		internal.ErrorSqlResponse(c, "failed to send msg")
		return
	}
	c.JSON(http.StatusOK, gin.H{"in": in, "conf": providerConfig})
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:provider", SendMsg)
}

package server

import (
	"accounts/global"
	"accounts/models"
	"accounts/msg/notify"
	"accounts/server/auth"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
	"strings"
)

// SendMsg godoc
//
//	@Summary	send message
//	@Schemes
//	@Description	send message
//	@Tags			msg
//	@Param			tenant		path	string			true	"tenant name"
//	@Param			providerId	path	integer			true	"provider id"
//	@Param			by			body	notify.SendInfo	true	"msg body"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{providerId} [post]
func SendMsg(c *gin.Context) {
	var in models.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	providerId := c.Param("providerId")
	var p models.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		internal.ErrorSqlResponse(c, "no such provider")
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	tenant := internal.GetTenant(c)
	authProvider, err := auth.GetAuthProvider(tenant.Id, p.Name)
	if err != nil {
		global.LOG.Error("get provider err: " + err.Error())
		internal.ErrorSqlResponse(c, "no such provider")
		return
	}

	var providerUser []string
	providerConfig := *authProvider.ProviderConfig()
	usersSlice := make([]string, len(in.Users))
	for i, v := range in.Users {
		usersSlice[i] = strings.Trim(v, "{}")
	}
	if err = global.DB.Table("provider_users pu").Select("pu.name").
		Joins("LEFT JOIN client_users as cu ON cu.user_id = pu.user_id").
		Where("pu.tenant_id = ? AND pu.provider_id = ? AND cu.sub IN (?)", tenant.Id, providerConfig["providerId"], usersSlice).
		Find(&providerUser).Error; err != nil {
		global.LOG.Error("get provider user err")
		internal.ErrorSqlResponse(c, "failed to get provider user")
		return
	}

	in.TenantId = tenant.Id

	var userDb []models.SendInfo
	for _, v := range in.Users {
		userDb = append(userDb, models.SendInfo{
			Link:       in.Link,
			UsersDB:    v,
			Sender:     in.Sender,
			Platform:   providerConfig["type"].(string),
			TenantId:   tenant.Id,
			Msg:        in.Msg,
			MsgType:    in.MsgType,
			Title:      in.Title,
			TitleColor: in.TitleColor,
			PngLink:    in.PngLink,
		})
	}

	// 调用InsertSendInfo函数插入数据到数据库
	if createErr := global.DB.Create(&userDb).Error; createErr != nil {
		global.LOG.Error("failed to insert SendInfo: " + createErr.Error())
		internal.ErrorSqlResponse(c, "failed to insert SendInfo")
		return
	}

	in.Users = providerUser
	in.Platform = providerConfig["type"].(string)
	if len(in.Users) == 0 {
		global.LOG.Warn("no provider user")
		internal.SuccessWithMessage(c, "no provider user")
		return
	}
	if err = notify.SendMsgToUsers(&in, providerConfig); err != nil {
		global.LOG.Error("send msg err: " + err.Error())
		internal.ErrorSqlResponse(c, "failed to send msg")
		return
	}
	internal.SuccessWithMessage(c, "ok")
}

// GetMsg godoc
//
//	@Summary	get message
//	@Schemes
//	@Description	get message
//	@Tags			msg
//	@Param			subId		path	integer			true	"sub id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{subId} [get]
func GetMsg(c *gin.Context) {
	subId := c.Param("subId")
	var SendInfo []models.SendInfo
	if err := internal.TenantDB(c).Model(&models.SendInfo{}).Where("users_db = ?", subId).Find(&SendInfo).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}
	var count int64
	if err := internal.TenantDB(c).Model(&models.SendInfo{}).Where("users_db = ?", subId).Count(&count).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}
	internal.SuccessWithDataAndTotal(c, SendInfo, count)
}

// MarkMsg godoc
//
//	@Summary	mark message read
//	@Schemes
//	@Description	mark message read
//	@Tags			msg
//	@Param			tenant		path	string			true	"tenant name"
//	@Param			msgId		path	integer			true	"msg id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{msgId} [put]
func MarkMsg(c *gin.Context) {
	var in models.SendInfo
	if err := c.ShouldBindUri(&in); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	var count int64
	if err := internal.TenantDB(c).Debug().Model(&models.SendInfo{}).Where("id = ?", in.Id).Count(&count).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to mark msg read")
		global.LOG.Error("mark msg read err: " + err.Error())
		return
	}
	if count == 0 {
		internal.SuccessWithMessage(c, "please check msg id")
		return
	}
	if err := internal.TenantDB(c).Model(&models.SendInfo{}).Where("id = ?", in.Id).Update("is_read", true).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to mark msg read")
		global.LOG.Error("mark msg read err: " + err.Error())
		return
	}
	internal.SuccessWithMessage(c, "mark msg read success")
}

// GetUnreadMsgCount godoc
//
//	@Summary	get unread message count
//	@Schemes
//	@Description	get unread message count
//	@Tags			msg
//	@Param			subId		path	integer			true	"sub id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{subId} [get]
func GetUnreadMsgCount(c *gin.Context) {
	subId := c.Param("subId")
	var count int64
	if err := internal.TenantDB(c).Debug().Model(&models.SendInfo{}).Where("users_db = ? AND is_read = ?", subId, false).Count(&count).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}
	internal.SuccessWithMessageAndData(c, "查询成功", count)
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:providerId", SendMsg)
	r.GET("/message/getMsg/:subId", GetMsg)
	r.PUT("/message/markMsg/:msgId", MarkMsg)
	r.GET("/message/unreadMsgCount/:subId", GetUnreadMsgCount)
}

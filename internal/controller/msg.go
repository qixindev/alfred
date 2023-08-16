package controller

import (
	"accounts/internal/controller/auth"
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/client/msg/notify"
	"accounts/pkg/global"
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
//	@Param			by			body	models.SendInfo	true	"msg body"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{providerId} [post]
func SendMsg(c *gin.Context) {
	var in model.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	providerId := c.Param("providerId")
	var p model.Provider
	if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get provider err")
		return
	}

	tenant := internal.GetTenant(c)
	authProvider, err := auth.GetAuthProvider(tenant.Id, p.Name)
	if err != nil {
		resp.ErrorUnknown(c, err, "get provider err")
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
		resp.ErrorSqlSelect(c, err, "get provider user err")
		return
	}

	in.TenantId = tenant.Id
	// 调用InsertSendInfo函数插入数据到数据库
	if createErr := global.DB.Create(&in).Error; createErr != nil {
		resp.ErrorSqlCreate(c, err, "failed to insert SendInfo")
		return
	}

	in.Users = providerUser
	in.Platform = providerConfig["type"].(string)
	if len(in.Users) == 0 {
		global.LOG.Warn("no provider user")
		resp.SuccessWithMessage(c, "no provider user")
		return
	}
	if err = notify.SendMsgToUsers(&in, providerConfig); err != nil {
		resp.ErrorUnknown(c, err, "failed to send msg")
		return
	}
	resp.SuccessWithMessage(c, "ok")
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
	var SendInfo []model.SendInfo
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("? = ANY(users)", subId).Find(&SendInfo).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get msg")
		return
	}
	var count int64
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("? = ANY(users)", subId).Count(&count).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get msg")
		return
	}
	resp.SuccessWithDataAndTotal(c, SendInfo, count)
}

// MarkMsg godoc
//
//	@Summary	mark message read
//	@Schemes
//	@Description	mark message read
//	@Tags			msg
//	@Param			subId		path	integer			true	"sub id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/MarkMsg [put]
func MarkMsg(c *gin.Context) {
	var in model.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("msg = ?", in.Msg).Update("is_read", in.IsRead).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "failed to mark msg read")
		return
	}
	resp.SuccessWithMessage(c, "mark msg read success")
}

// GetUnreadMsgCount godoc
//
//	@Summary	get unread message count
//	@Schemes
//	@Description	get unread message count
//	@Tags			msg
//	@Param			subId		path	integer			true	"sub id"
//	@Success		200
//	@Router			/accounts/{tenant}/message/unreadMsgCount/{subId} [get]
func GetUnreadMsgCount(c *gin.Context) {
	subId := c.Param("subId")
	var count int64
	if err := internal.TenantDB(c).Debug().Model(&model.SendInfo{}).Where("? = ANY(users)", subId).Where("is_read = ?", false).Count(&count).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get unread msg count")
		return
	}
	resp.SuccessWithMessageAndData(c, "查询成功", count)
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:providerId", SendMsg)
	r.GET("/message/:subId", GetMsg)
	r.PUT("/message/markMsgRead", MarkMsg)
	r.GET("/message/unreadMsgCount/:subId", GetUnreadMsgCount)
}

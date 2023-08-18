package controller

import (
	"accounts/internal/controller/auth"
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/client/msg/notify"
	"accounts/pkg/global"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		resp.ErrorSqlResponse(c, "no such provider")
		global.LOG.Error("get provider err: " + err.Error())
		return
	}

	tenant := internal.GetTenant(c)
	in.TenantId = tenant.Id
	authProvider, err := auth.GetAuthProvider(tenant.Id, p.Name)
	if err != nil {
		global.LOG.Error("get provider err: " + err.Error())
		resp.ErrorSqlResponse(c, "no such provider")
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
		resp.ErrorSqlResponse(c, "failed to get provider user")
		return
	}

	var userDb []model.SendInfo
	for _, v := range in.Users {
		userDb = append(userDb, model.SendInfo{
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
		resp.ErrorSqlResponse(c, "failed to insert SendInfo")
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
		global.LOG.Error("send msg err: " + err.Error())
		resp.ErrorSqlResponse(c, "failed to send msg")
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
//	@Param			page		query	integer			false	"pageNum"
//	@Param			pageSize	query	integer			false	"pageSize"
//	@Success		200
//	@Router			/accounts/{tenant}/message/{subId} [get]
func GetMsg(c *gin.Context) {
	subId := c.Param("subId")
	var SendInfo []model.SendInfo
	var SendInfoDB []model.SendInfoDB

	// 获取页码，默认为1
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	// 获取每页显示的数据数量，默认为10
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 通过JOIN查询获取Message数据和发送者、接收者的显示名
	if err := global.DB.Debug().
		Table("message").Debug().
		Select("message.*, u1.display_name as sender_name, u2.display_name as receiver_name").
		Joins("LEFT JOIN client_users cu1 ON message.sender = cu1.sub").
		Joins("LEFT JOIN client_users cu2 ON message.users_db = cu2.sub").
		Joins("LEFT JOIN users u1 ON cu1.user_id = u1.id").
		Joins("LEFT JOIN users u2 ON cu2.user_id = u2.id").
		Where("message.users_db = ?", subId).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&SendInfoDB).Error; err != nil {
		resp.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}

	// 获取消息
	var total int64
	if err := internal.TenantDB(c).Debug().Model(&model.SendInfo{}).
		Where("users_db = ?", subId).Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&SendInfo).Count(&total).Error; err != nil {
		resp.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}

	for i, v := range SendInfo {
		for _, v2 := range SendInfoDB {
			if v.Id == v2.Id {
				SendInfo[i].SenderName = v2.SenderName
				SendInfo[i].ReceiverName = v2.ReceiverName
			}
		}
	}

	resp.SuccessWithDataAndTotal(c, SendInfo, total)
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
	var in model.SendInfo
	if err := c.ShouldBindUri(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	var count int64
	if err := internal.TenantDB(c).Debug().Model(&model.SendInfo{}).Where("id = ?", in.Id).Count(&count).Error; err != nil {
		resp.ErrorSqlResponse(c, "failed to mark msg read")
		global.LOG.Error("mark msg read err: " + err.Error())
		return
	}
	if count == 0 {
		resp.SuccessWithMessage(c, "please check msg id")
		return
	}
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("id = ?", in.Id).Update("is_read", true).Error; err != nil {
		resp.ErrorSqlResponse(c, "failed to mark msg read")
		global.LOG.Error("mark msg read err: " + err.Error())
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
//	@Router			/accounts/{tenant}/unreadMsgCount/{subId} [get]
func GetUnreadMsgCount(c *gin.Context) {
	subId := c.Param("subId")
	var count int64
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("users_db = ? AND is_read = ?", subId, false).Count(&count).Error; err != nil {
		resp.ErrorSqlResponse(c, "failed to get msg")
		global.LOG.Error("get msg err: " + err.Error())
		return
	}
	resp.SuccessWithMessageAndData(c, "查询成功", count)
}

func AddMsgRouter(r *gin.RouterGroup) {
	r.POST("/message/:providerId", SendMsg)
	r.GET("/message/getMsg/:subId", GetMsg)
	r.PUT("/message/markMsg/:msgId", MarkMsg)
	r.GET("/message/unreadMsgCount/:subId", GetUnreadMsgCount)
}

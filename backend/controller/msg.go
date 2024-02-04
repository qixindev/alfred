package controller

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/client/msg/notify"
	"alfred/backend/pkg/config/env"
	"alfred/backend/pkg/global"
	"alfred/backend/service/auth"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SendMsg
// @Summary	send message
// @Tags	msg
// @Param	tenant	path		string	true	"tenant"	default(default)
// @Param	providerId	path	integer			true	"provider id"
// @Param	by			body	model.SendInfo	true	"msg body"
// @Success	200
// @Router	/accounts/{tenant}/message/{providerId} [post]
func SendMsg(c *gin.Context) {
	var in model.SendInfo
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	var providerUser []string
	var providerConfig gin.H

	tenant := internal.GetTenant(c)
	in.TenantId = tenant.Id
	providerId := c.Param("providerId")

	if providerId == "0" { // 站内消息
		providerConfig = make(gin.H)
		providerConfig["type"] = env.PlatformAlfred
	} else {
		var p model.Provider
		if err := internal.TenantDB(c).First(&p, "id = ?", providerId).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "get provider err")
			return
		}

		_, authProvider, err := auth.GetAuthProvider(tenant.Id, p.Name)
		if err != nil {
			resp.ErrorUnknown(c, err, "get provider err")
			return
		}

		providerConfig = *authProvider.ProviderConfig()
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
			SendAt:     time.Now(),
		})
	}

	// 调用InsertSendInfo函数插入数据到数据库
	if err := global.DB.Create(&userDb).Error; err != nil {
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

	if err := notify.SendMsgToUsers(&in, providerConfig); err != nil {
		resp.ErrorUnknown(c, err, "failed to send msg")
		return
	}
	resp.SuccessWithMessage(c, "ok")
}

// GetMsg
// @Summary	get message
// @Tags	msg
// @Param	subId	path	string	true	"sub id"
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	msgType	query	string	false	"msg type"
// @Param	page	query	integer	false	"pageNum"
// @Param	pageSize	query	integer	false	"pageSize"
// @Success	200
// @Router	/accounts/{tenant}/message/getMsg/{subId} [get]
func GetMsg(c *gin.Context) {
	subId := c.Param("subId")
	var SendInfo []model.SendInfo
	var SendInfoDB []model.SendInfoDB

	tenant := internal.GetTenant(c)

	// 获取页码，默认为1
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	// 获取每页显示的数据数量，默认为10
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 获取消息类型参数
	messageType := c.Query("msgType") //// 获取消息

	var total int64
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).
		Where("users_db = ?", subId).
		Find(&SendInfo).Count(&total).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get msg")
		return
	}

	query := global.DB.Table("message").
		Select("message.*, u1.display_name as sender_name, u2.display_name as receiver_name, u1.avatar").
		Joins("LEFT JOIN client_users cu1 ON message.sender = cu1.sub").
		Joins("LEFT JOIN client_users cu2 ON message.users_db = cu2.sub").
		Joins("LEFT JOIN users u1 ON cu1.user_id = u1.id").
		Joins("LEFT JOIN users u2 ON cu2.user_id = u2.id").
		Where("message.users_db = ? AND message.tenant_id = ?", subId, tenant.Id)

	if messageType != "" {
		query = query.Where("message.msg_type = ?", messageType)
	}

	if err := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("send_at desc").
		Find(&SendInfoDB).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get msg")
		return
	}

	for i, v := range SendInfoDB {
		for _, v2 := range SendInfo {
			if v.Id == v2.Id {
				SendInfo[i].SenderName = v2.SenderName
				SendInfo[i].ReceiverName = v2.ReceiverName
			}
		}
	}

	resp.SuccessWithDataAndTotal(c, SendInfoDB, total)
}

// MarkMsg
// @Summary	mark message read
// @Tags	msg
// @Param	tenant	path		string	true	"tenant"	default(default)
// @Param	msgId	path	integer	true	"msg id"
// @Success	200
// @Router	/accounts/{tenant}/message/markMsg/{msgId} [put]
func MarkMsg(c *gin.Context) {
	var in model.SendInfo
	if err := c.ShouldBindUri(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	var count int64
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("id = ?", in.Id).Count(&count).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to mark msg read")
		return
	}
	if count == 0 {
		resp.SuccessWithMessage(c, "please check msg id")
		return
	}
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("id = ?", in.Id).Update("is_read", true).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "failed to mark msg read")
		return
	}
	resp.SuccessWithMessage(c, "mark msg read success")
}

// GetUnreadMsgCount
// @Summary	get unread message count
// @Tags	msg
// @Param	subId	path	string	true	"sub id"
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200
// @Router	/accounts/{tenant}/unreadMsgCount/{subId} [get]
func GetUnreadMsgCount(c *gin.Context) {
	subId := c.Param("subId")
	var count int64
	if err := internal.TenantDB(c).Model(&model.SendInfo{}).Where("users_db = ? AND is_read = ?", subId, false).Count(&count).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "failed to get msg")
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

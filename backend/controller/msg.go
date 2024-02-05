package controller

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/client/msg/notify"
	"alfred/backend/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
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

	msgService := service.NewMsgService()
	providerId := c.Param("providerId")
	msgService.Tenant = internal.GetTenant(c)
	in.TenantId = msgService.Tenant.Id

	providerConfig, err := msgService.ProcessMsg(providerId, in)
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to process msg")
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
// @Param	subId		path	string	true	"sub id"
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	msgTypes	query	string	false	"msg type"
// @Param	page		query	integer	false	"pageNum"
// @Param	pageSize	query	integer	false	"pageSize"
// @Success	200
// @Router	/accounts/{tenant}/message/getMsg/{subId} [get]
func GetMsg(c *gin.Context) {
	msgService := service.NewMsgService()
	msgService.Tenant = internal.GetTenant(c)

	subId := c.Param("subId")
	msgTypes := strings.Split(c.Query("msgTypes"), ",")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	list, total, err := msgService.GetMsgList(subId, msgTypes, pageNum, pageSize)
	if err != nil {
		return
	}
	resp.SuccessWithDataAndTotal(c, list, total)
}

// MarkMsg
// @Summary	mark message read
// @Tags	msg
// @Param	tenant	path		string	true	"tenant"	default(default)
// @Param	msgId	path	integer	true	"msg id"
// @Success	200
// @Router	/accounts/{tenant}/message/markMsg/{msgId} [put]
func MarkMsg(c *gin.Context) {
	msgService := service.NewMsgService()
	msgService.Tenant = internal.GetTenant(c)

	var in model.SendInfo
	if err := c.ShouldBindUri(&in); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	if err := msgService.MarkMsg(in.Id, true); err != nil {
		resp.ErrorUnknown(c, err, "failed to mark msg read")
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
	msgService := service.NewMsgService()
	msgService.Tenant = internal.GetTenant(c)

	subId := c.Param("subId")
	count, err := msgService.GetMsgCount(subId, false)
	if err != nil {
		resp.ErrorUnknown(c, err, "failed to get unread msg count")
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

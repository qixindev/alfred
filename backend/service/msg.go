package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/config/env"
	"alfred/backend/pkg/global"
	"alfred/backend/service/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// MsgService 包含消息相关的服务逻辑
type MsgService struct{}

// NewMsgService 创建一个新的消息服务
func NewMsgService() *MsgService {
	return &MsgService{}
}

func (*MsgService) ProcessMsg(providerId string, in model.SendInfo, tenant *model.Tenant) (gin.H, error) {
	var providerUser []string
	var providerConfig gin.H

	if providerId == "0" { // 站内消息
		providerConfig = make(gin.H)
		providerConfig["type"] = env.PlatformAlfred
	} else {
		var p model.Provider
		if err := global.DB.First(&p, "id = ? AND tenant_id = ?", providerId, tenant.Id).Error; err != nil {
			return gin.H{}, err
		}

		_, authProvider, err := auth.GetAuthProvider(tenant.Id, p.Name)
		if err != nil {
			return gin.H{}, err
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
			return gin.H{}, err
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
		return gin.H{}, err
	}

	in.Users = providerUser
	in.Platform = providerConfig["type"].(string)
	if len(in.Users) == 0 {
		global.LOG.Warn("no provider user")
		return providerConfig, errors.New("no provider user")
	}
	return providerConfig, nil
}

func (*MsgService) GetMsgList(subId string, tenant *model.Tenant, msgTypes []string, pageNum int, pageSize int, SendInfoDB []model.SendInfoDB, SendInfo []model.SendInfo) ([]model.SendInfoDB, int64, error) {
	query := global.DB.Table("message").
		Select("message.*, u1.display_name as sender_name, u2.display_name as receiver_name, u1.avatar").
		Joins("LEFT JOIN client_users cu1 ON message.sender = cu1.sub").
		Joins("LEFT JOIN client_users cu2 ON message.users_db = cu2.sub").
		Joins("LEFT JOIN users u1 ON cu1.user_id = u1.id").
		Joins("LEFT JOIN users u2 ON cu2.user_id = u2.id").
		Where("message.users_db = ? AND message.tenant_id = ?", subId, tenant.Id)

	if len(msgTypes) > 0 && msgTypes[0] != "" {
		query = query.Where("message.msg_type IN ?", msgTypes)
	}

	if err := query.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("send_at desc").
		Find(&SendInfoDB).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := global.DB.Model(&model.SendInfo{}).
		Where("users_db = ? AND msg_type IN (?) AND tenant_id = ?", subId, msgTypes, tenant.Id).
		Find(&SendInfo).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	for i, v := range SendInfoDB {
		for _, v2 := range SendInfo {
			if v.Id == v2.Id {
				SendInfo[i].SenderName = v2.SenderName
				SendInfo[i].ReceiverName = v2.ReceiverName
			}
		}
	}
	return SendInfoDB, total, nil
}

func (*MsgService) MarkMsgAsRead(msgId uint, tenant *model.Tenant) error {
	var count int64
	if err := global.DB.Model(&model.SendInfo{}).Where("id = ? AND tenant_id = ?", msgId, tenant.Id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("msg not found")
	}
	if err := global.DB.Model(&model.SendInfo{}).Where("id = ? AND tenant_id = ?", msgId, tenant.Id).Update("is_read", true).Error; err != nil {
		return err
	}
	return nil
}

func (*MsgService) GetUnreadMsgCount(subId string, tenant *model.Tenant) (int64, error) {
	var count int64
	if err := global.DB.Model(&model.SendInfo{}).Where("users_db = ? AND is_read = ? AND tenant_id = ?", subId, false, tenant.Id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

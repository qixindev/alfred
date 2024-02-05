package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
)

func MarkMsgAsRead(msgId uint, tenant *model.Tenant) error {
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

func GetUnreadMsgCount(subId string, tenant *model.Tenant) (int64, error) {
	var count int64
	if err := global.DB.Model(&model.SendInfo{}).Where("users_db = ? AND is_read = ? AND tenant_id = ?", subId, false, tenant.Id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

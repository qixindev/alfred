package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
)

func DeleteDevice(tenantId uint, deviceId string) error {
	if err := global.DB.Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).
		Delete(model.GroupDevice{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).
		Delete(model.DeviceSecret{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, deviceId).
		Delete(model.Device{}).Error; err != nil {
		return err
	}
	return nil
}

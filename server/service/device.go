package service

import (
	"accounts/global"
	"accounts/models"
)

func DeleteDevice(tenantId uint, deviceId string) error {
	if err := global.DB.Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).
		Delete(models.GroupDevice{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).
		Delete(models.DeviceSecret{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, deviceId).
		Delete(models.Device{}).Error; err != nil {
		return err
	}
	return nil
}

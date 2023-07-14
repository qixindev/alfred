package service

import (
	"accounts/global"
	"accounts/models"
)

func DeleteDevice(tenantId uint, deviceId string) error {
	if err := global.DB.Delete(models.GroupDevice{}).
		Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).Error; err != nil {
		return err
	}

	if err := global.DB.Delete(models.DeviceCode{}).
		Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).Error; err != nil {
		return err
	}

	if err := global.DB.Delete(models.DeviceSecret{}).
		Where("tenant_id = ? AND device_id = ?", tenantId, deviceId).Error; err != nil {
		return err
	}

	if err := global.DB.Delete(models.Device{}).
		Where("tenant_id = ? AND id = ?", tenantId, deviceId).Error; err != nil {
		return err
	}
	return nil
}

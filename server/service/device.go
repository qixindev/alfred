package service

import (
	"accounts/global"
	"accounts/models"
	"github.com/google/uuid"
)

func GetDeviceAndSecret(deviceCode string, tenantId uint) (string, string, error) {
	device := models.Device{
		Name:     deviceCode,
		TenantId: tenantId,
	}
	if err := global.DB.Create(&device).Error; err != nil {
		return "", "", err
	}

	deviceSecret := models.DeviceSecret{
		Name:     "default",
		DeviceId: device.Id,
		Secret:   uuid.NewString(),
	}
	if err := global.DB.Create(&deviceSecret).Error; err != nil {
		return "", "", err
	}

	return device.Name, deviceSecret.Secret, nil
}

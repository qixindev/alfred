package service

import (
	"accounts/internal/global"
	"accounts/pkg/models"
	"time"
)

func ClearDeviceCode(code string) {
	earliest := time.Now().Add(-2 * time.Minute)
	if err := global.DB.Delete(&models.DeviceCode{}, "user_code = ? OR created_at < ?", code, earliest).Error; err != nil {
		global.LOG.Error("delete device code err: " + err.Error())
	}
}

func ClearTokenCode(code string) {
	earliest := time.Now().Add(-2 * time.Minute)
	if err := global.DB.Delete(&models.TokenCode{}, "code = ? OR created_at < ?", code, earliest).Error; err != nil {
		global.LOG.Error("delete token code err: " + err.Error())
	}
}

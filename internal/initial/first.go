package initial

import (
	"accounts/internal/global"
	"accounts/pkg/config/env"
	"accounts/pkg/models"
	"accounts/utils"
)

const (
	DefaultTenant = "default"
	DefaultClient = "default"
	DefaultUser   = "admin"
	DefaultPwd    = "admin"
)

func CheckFirstRun() error {
	if env.GetReleaseType() == "first" {
		var tenant models.Tenant
		if err := global.DB.First(&tenant, "name = ?", DefaultTenant).Error; err != nil {
			return initFirstRun()
		}
		if err := migrateDB(); err != nil {
			return err
		}
	}

	return nil
}

func initFirstRun() error {
	if err := migrateDB(); err != nil {
		return err
	}

	if err := insertDB(); err != nil {
		return err
	}

	if _, err := utils.LoadRsaPublicKeys(DefaultTenant); err != nil {
		return err
	}

	return nil
}

func insertDB() error {
	var tenant models.Tenant
	tenant.Name = DefaultTenant
	if err := global.DB.Create(&tenant).Error; err != nil {
		return err
	}

	if err := global.DB.Create(&models.Client{Id: DefaultClient, Name: DefaultClient, TenantId: tenant.Id}).Error; err != nil {
		return err
	}

	adminPwd, err := utils.HashPassword(DefaultPwd)
	if err != nil {
		return err
	}
	if err = global.DB.Create(&models.User{
		Username:         DefaultUser,
		PasswordHash:     adminPwd,
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		Disabled:         false,
		TenantId:         tenant.Id,
		Role:             "admin",
	}).Error; err != nil {
		return err
	}

	return nil
}

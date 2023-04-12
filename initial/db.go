package initial

import (
	"accounts/config/env"
	"accounts/global"
	"accounts/models"
	"accounts/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CheckFirstRun() error {
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "name = ?", "default1").Error; err != nil {
		return insertDB()
	}

	return nil
}

func InitDB() error {
	dsn := global.CONFIG.Pgsql.ConfigDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	global.DB = db
	if env.GetReleaseType() == "first" {
		if err = migrateDB(); err != nil {
			return err
		}

		if err = CheckFirstRun(); err != nil {
			return err
		}
	}

	return nil
}

func insertDB() error {
	var tenant models.Tenant
	tenant.Name = "default"
	if err := global.DB.Create(&tenant).Error; err != nil {
		return err
	}

	if err := global.DB.Create(&models.Client{Id: "default", Name: "default", TenantId: tenant.Id}).Error; err != nil {
		return err
	}

	adminPwd, err := utils.HashPassword("admin")
	if err != nil {
		return err
	}
	if err = global.DB.Create(&models.User{
		Username:         "admin",
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

func migrateDB() error {
	migrateList := []any{
		&models.Tenant{},
		&models.User{},
		&models.Group{},
		&models.Client{},
		&models.ClientUser{},
		&models.Device{},
		&models.DeviceSecret{},
		&models.DeviceCode{},
		&models.GroupUser{},
		&models.GroupDevice{},
		&models.RedirectUri{},
		&models.ClientSecret{},
		&models.TokenCode{},
		&models.ProviderUser{},
		&models.Provider{},
		&models.ProviderOAuth2{},
		&models.ProviderDingTalk{},
		&models.ProviderWeCom{},
		&models.SmsConnector{},
		&models.SmsTcloud{},
		&models.ResourceType{},
		&models.ResourceTypeAction{},
		&models.ResourceTypeRole{},
		&models.ResourceTypeRoleAction{},
		&models.Resource{},
		&models.ResourceRoleUser{},
	}

	if err := global.DB.AutoMigrate(migrateList...); err != nil {
		fmt.Println("migrate db err: ", err)
		return err
	}
	return nil
}

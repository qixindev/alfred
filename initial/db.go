package initial

import (
	"accounts/global"
	"accounts/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() error {
	dsn := global.CONFIG.Pgsql.ConfigDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	global.DB = db
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

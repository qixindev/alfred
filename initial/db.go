package initial

import (
	"accounts/internal/model"
	"accounts/pkg/global"
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
		&model.Tenant{},
		&model.User{},
		&model.Group{},
		&model.Client{},
		&model.ClientUser{},
		&model.Device{},
		&model.DeviceSecret{},
		&model.DeviceCode{},
		&model.GroupUser{},
		&model.GroupDevice{},
		&model.RedirectUri{},
		&model.ClientSecret{},
		&model.TokenCode{},
		&model.ProviderUser{},
		&model.Provider{},
		&model.ProviderOAuth2{},
		&model.ProviderDingTalk{},
		&model.ProviderWeCom{},
		&model.SmsConnector{},
		&model.SmsTcloud{},
		&model.ResourceType{},
		&model.ResourceTypeAction{},
		&model.ResourceTypeRole{},
		&model.ResourceTypeRoleAction{},
		&model.Resource{},
		&model.ResourceRoleUser{},
		&model.SendInfo{},
		&model.PhoneVerification{},
	}

	if err := global.DB.AutoMigrate(migrateList...); err != nil {
		fmt.Println("migrate db err: ", err)
		return err
	}
	return nil
}

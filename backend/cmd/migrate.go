package cmd

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"fmt"
	"os"
)

func getMigrateModel() []any {
	return []any{
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
		&model.ProviderWechat{},
		&model.ProviderSms{},
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
}

func migrateDB() {
	if err := initSystem(); err != nil {
		fmt.Println("init system err:", err)
		os.Exit(1)
		return
	}
	migrateList := getMigrateModel()
	if err := global.DB.AutoMigrate(migrateList...); err != nil {
		fmt.Println("[Error] migrate db err: ", err)
		os.Exit(2)
		return
	}

	fmt.Println("===== Success =====")
}

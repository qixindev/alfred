package cmd

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"fmt"
	"os"
)

func migrateDB() {
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
		os.Exit(2)
		return
	}
	fmt.Println("===== Success =====")
}

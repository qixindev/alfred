package service

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"
)

func DeleteTenant(tenant model.Tenant) error {
	if tenant.Id == 0 {
		return errors.New("delete invalid tenant")
	}
	delList := []any{
		model.RedirectUri{}, model.TokenCode{},
		model.ResourceTypeRoleAction{}, model.ResourceRoleUser{},
		model.ResourceTypeAction{}, model.Resource{}, model.ResourceTypeRole{}, model.ResourceType{},

		model.ProviderUser{}, model.ProviderDingTalk{}, model.ProviderWeCom{},
		model.ProviderSms{}, model.ProviderOAuth2{}, model.Provider{},

		model.GroupUser{}, model.GroupDevice{}, model.Group{},
		model.DeviceSecret{}, model.DeviceCode{}, model.Device{},
		model.ClientUser{}, model.ClientSecret{}, model.Client{},
		model.User{},
	}

	for _, v := range delList {
		if err := global.DB.Model(v).Where("tenant_id = ?", tenant.Id).Delete(v).Error; err != nil {
			return err
		}
	}
	if err := global.DB.Where("id = ?", tenant.Id).Delete(model.Tenant{}).Error; err != nil {
		return err
	}

	return nil
}

package service

import (
	"accounts/internal/model"
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
		model.ProviderOAuth2{}, model.Provider{},

		model.GroupUser{}, model.GroupDevice{}, model.Group{},
		model.DeviceSecret{}, model.DeviceCode{}, model.Device{},
		model.ClientUser{}, model.ClientSecret{}, model.Client{},
		model.User{},
	}

	if err := deleteSource(model.Tenant{}, delList, tenant.Id, "tenant_id"); err != nil {
		return err
	}
	return nil
}

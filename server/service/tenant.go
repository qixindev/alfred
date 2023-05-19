package service

import (
	"accounts/models"
	"errors"
)

func DeleteTenant(tenant models.Tenant) error {
	if tenant.Id == 0 {
		return errors.New("delete invalid tenant")
	}
	delList := []any{
		models.RedirectUri{}, models.TokenCode{},
		models.ResourceTypeRoleAction{}, models.ResourceRoleUser{},
		models.ResourceTypeAction{}, models.Resource{}, models.ResourceTypeRole{}, models.ResourceType{},

		models.ProviderUser{}, models.ProviderDingTalk{}, models.ProviderWeCom{},
		models.ProviderOAuth2{}, models.Provider{},

		models.GroupUser{}, models.GroupDevice{}, models.Group{},
		models.DeviceSecret{}, models.DeviceCode{}, models.Device{},
		models.ClientUser{}, models.ClientSecret{}, models.Client{},
		models.User{},
	}

	if err := deleteSource(models.Tenant{}, delList, tenant.Id, "tenant_id"); err != nil {
		return err
	}
	return nil
}

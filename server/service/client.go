package service

import (
	"accounts/models"
	"errors"
)

func DeleteClient(clientId string) error {
	if clientId == "" {
		return errors.New("delete invalid client")
	}
	delList := []any{
		models.RedirectUri{}, models.TokenCode{},
		models.ResourceTypeRoleAction{}, models.ResourceRoleUser{},
		models.ResourceTypeAction{}, models.Resource{}, models.ResourceTypeRole{}, models.ResourceType{},

		models.ProviderUser{}, models.ProviderDingTalk{}, models.ProviderWeCom{},
		models.ProviderOAuth2{}, models.Provider{},

		models.GroupUser{}, models.GroupDevice{}, models.Group{},
		models.DeviceSecret{}, models.DeviceCode{}, models.Device{},
		models.ClientUser{}, models.ClientSecret{},
	}

	if err := deleteSource(models.Client{}, delList, clientId, "client_id"); err != nil {
		return err
	}
	return nil
}

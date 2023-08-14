package service

import (
	"accounts/internal/global"
	"accounts/pkg/models"
	"errors"
	"net/url"
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

func IsValidateUri(tenantId uint, clientId, uri string) error {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return err
	}

	host := parsedURL.Scheme + "://" + parsedURL.Host
	if err = global.DB.First(&models.RedirectUri{}, "tenant_id = ? AND client_id = ? AND redirect_uri = ?", tenantId, clientId, host).Error; err != nil {
		return err
	}

	return nil
}

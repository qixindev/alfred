package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
	"net/url"
)

func DeleteClient(tenantId uint, clientId string) error {
	if clientId == "" {
		return errors.New("delete invalid client")
	}
	delList := []any{
		model.RedirectUri{}, model.TokenCode{}, model.ResourceType{},
		model.ClientUser{}, model.ClientSecret{},
	}

	if err := deleteSource(tenantId, delList, clientId, "client_id"); err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, clientId).Delete(model.Client{}).Error; err != nil {
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
	if err = global.DB.First(&model.RedirectUri{}, "tenant_id = ? AND client_id = ? AND redirect_uri = ?", tenantId, clientId, host).Error; err != nil {
		return err
	}

	return nil
}

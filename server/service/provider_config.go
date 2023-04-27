package service

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/types"
	"errors"
)

func GetProviderModel(t string) (models.ItfProvider, error) {
	switch t {
	case "oauth2":
		return &models.ProviderOAuth2{}, nil
	case "dingtalk":
		return &models.ProviderDingTalk{}, nil
	case "wecom":
		return &models.ProviderWeCom{}, nil
	}
	return nil, errors.New("no such type")
}

func CreateProviderConfig(p types.ReqProvider) error {
	it, err := GetProviderModel(p.Type)
	if err != nil {
		return err
	}

	provider := models.Provider{Name: p.Name, Type: p.Type, TenantId: p.TenantId}
	if err = global.DB.Create(&provider).Error; err != nil {
		return err
	}

	p.ProviderId = provider.Id
	if err = global.DB.Create(it.Save(p)).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProviderConfig(p types.ReqProvider) error {
	it, err := GetProviderModel(p.Type)
	if err != nil {
		return err
	}

	provider := models.Provider{Name: p.Name, Type: p.Type}
	if err = global.DB.Where("tenant_id = ? AND id = ? AND type = ?", p.TenantId, p.ProviderId, p.Type).
		Updates(&provider).Error; err != nil {
		return err
	}

	if err = global.DB.Where("tenant_id = ? AND provider_id = ?", p.TenantId, p.ProviderId).
		Updates(it.Save(p)).Error; err != nil {
		return err
	}

	return nil
}

func GetProvider(tenantId uint, providerId uint, t string) (any, error) {
	pr, err := GetProviderModel(t)
	if err != nil {
		return nil, err
	}

	if err = global.DB.Model(pr).Where("tenant_id = ? AND provider_id = ?", tenantId, providerId).
		Preload("Provider").First(pr).Error; err != nil {
		return nil, err
	}

	return pr.Dto(), nil
}

func DeleteProviderConfig(p models.Provider) error {
	pr, err := GetProviderModel(p.Type)
	if err != nil {
		return err
	}
	if err = global.DB.Where("tenant_id = ? AND provider_id = ?", p.TenantId, p.Id).
		Delete(pr).Error; err != nil {
		return err
	}

	return nil
}

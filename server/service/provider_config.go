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

func GetProviderUsers(tenantId uint, providerId uint, client string) (any, error) {
	var users []models.ProviderUser
	if err := global.DB.Table("provider_users as pu").Select("cu.sub", "u.display_name", "pu.provider_id").
		Joins("LEFT JOIN users u ON u.id = pu.user_id").
		Joins("LEFT JOIN client_users cu ON cu.user_id = pu.user_id").
		Where("pu.provider_id = ? AND pu.tenant_id = ? AND cu.client_id = ?", providerId, tenantId, client).
		Preload("Provider").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
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

package service

import (
	"accounts/internal/endpoint/req"
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"
)

func GetProviderModel(t string) (model.ItfProvider, error) {
	switch t {
	case "oauth2":
		return &model.ProviderOAuth2{}, nil
	case "dingtalk":
		return &model.ProviderDingTalk{}, nil
	case "wecom":
		return &model.ProviderWeCom{}, nil
	}
	return nil, errors.New("no such type")
}

func CreateProviderConfig(p req.ReqProvider) error {
	it, err := GetProviderModel(p.Type)
	if err != nil {
		return err
	}

	provider := model.Provider{Name: p.Name, Type: p.Type, TenantId: p.TenantId}
	if err = global.DB.Create(&provider).Error; err != nil {
		return err
	}

	p.ProviderId = provider.Id
	if err = global.DB.Create(it.Save(p)).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProviderConfig(p req.ReqProvider) error {
	it, err := GetProviderModel(p.Type)
	if err != nil {
		return err
	}

	provider := model.Provider{Name: p.Name, Type: p.Type}
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

func GetProviderUsers(tenantId uint, providerId uint) (any, error) {
	var users []model.ProviderUser
	if err := global.DB.Table("provider_users as pu").
		Select("pu.provider_id", "cu.sub", "u.display_name").
		Joins("LEFT JOIN users u ON u.id = pu.user_id").
		Joins("LEFT JOIN client_users cu ON cu.user_id = pu.user_id").
		Where("pu.provider_id = ? AND pu.tenant_id = ?", providerId, tenantId).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteProviderConfig(p model.Provider) error {
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

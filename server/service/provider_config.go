package service

import (
	"accounts/global"
	"accounts/models"
	"errors"
	"gorm.io/gorm"
)

func CreateProviderConfig(p models.Provider) error {
	var err error
	switch p.Type {
	case "oauth2":
		err = global.DB.Create(&models.ProviderOAuth2{
			ProviderId:   p.Id,
			TenantId:     p.TenantId,
			ClientId:     p.ClientId,
			ClientSecret: p.ClientSecret,
		}).Error
	case "dingtalk":
		err = global.DB.Create(&models.ProviderDingTalk{
			ProviderId: p.Id,
			TenantId:   p.TenantId,
			AgentId:    p.AgentId,
			AppKey:     p.ClientId,
			AppSecret:  p.ClientSecret,
		}).Error
	case "wecom":
		err = global.DB.Create(&models.ProviderWeCom{
			ProviderId: p.Id,
			TenantId:   p.TenantId,
			CorpId:     p.AgentId,
			AgentId:    p.ClientId,
			AppSecret:  p.ClientSecret,
		}).Error
	default:
		return errors.New("no such provider type")
	}
	return err
}

func UpdateProviderConfig(p models.Provider) error {
	tx := global.DB.Debug().Where("tenant_id = ? AND provider_id = ?", p.TenantId, p.Id)
	switch p.Type {
	case "oauth2":
		tx.Updates(&models.ProviderOAuth2{
			ProviderId:   p.Id,
			TenantId:     p.TenantId,
			ClientId:     p.ClientId,
			ClientSecret: p.ClientSecret,
		})
	case "dingtalk":
		tx.Updates(&models.ProviderDingTalk{
			ProviderId: p.Id,
			TenantId:   p.TenantId,
			AgentId:    p.AgentId,
			AppKey:     p.ClientId,
			AppSecret:  p.ClientSecret,
		})
	case "wecom":
		tx.Updates(&models.ProviderWeCom{
			ProviderId: p.Id,
			TenantId:   p.TenantId,
			CorpId:     p.AgentId,
			AgentId:    p.ClientId,
			AppSecret:  p.ClientSecret,
		})
	default:
		return errors.New("no such provider type")
	}

	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func GetProvider(tenantId uint, providerId uint, t string) (map[string]any, error) {
	var tx *gorm.DB
	var res map[string]any
	switch t {
	case "oauth2":
		tx = global.DB.Model(models.ProviderOAuth2{})
	case "dingtalk":
		tx = global.DB.Model(models.ProviderDingTalk{})
	case "wecom":
		tx = global.DB.Model(models.ProviderWeCom{})
	default:
		return nil, errors.New("no such provider type")
	}

	if err := tx.Where("tenant_id = ? AND provider_id = ?", tenantId, providerId).
		Preload("Provider").First(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

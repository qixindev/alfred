package service

import (
	"accounts/global"
	"accounts/models"
	"errors"
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

func GetProvider(tenantId uint, providerId uint, t string) (any, error) {
	tx := global.DB.Where("tenant_id = ? AND provider_id = ?", tenantId, providerId)
	var err error
	oauth2, ding, wecom := models.ProviderOAuth2{}, models.ProviderDingTalk{}, models.ProviderWeCom{}
	switch t {
	case "oauth2":
		err = tx.Model(oauth2).Preload("Provider").First(&oauth2).Error
		return oauth2, err
	case "dingtalk":
		err = tx.Model(ding).Preload("Provider").First(&ding).Error
		return ding, err
	case "wecom":
		err = tx.Model(wecom).Preload("Provider").First(&wecom).Error
		return wecom, err
	}

	return nil, errors.New("no such provider type")
}

func IsValidType(t string) bool {
	return t == "oauth2" || t == "dingtalk" || t == "wecom"
}

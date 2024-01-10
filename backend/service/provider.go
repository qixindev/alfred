package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
)

func DeleteProvider(tenantId uint, id uint) error {
	delList := []any{
		model.ProviderDingTalk{},
		model.ProviderWeCom{},
		model.ProviderOAuth2{},
	}
	if err := deleteSource(tenantId, delList, id, "provider_id"); err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, id).Delete(model.Provider{}).Error; err != nil {
		return err
	}

	return nil
}

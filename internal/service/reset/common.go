package reset

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"
)

type ProviderSms interface {
	ResetAuth(string, uint) (string, error)
}

func GetResetAuthProvider(tenantId uint, providerName string) (*ProviderResetSms, error) {
	var provider model.Provider
	if err := global.DB.First(&provider, "tenant_id = ? AND name = ?", tenantId, providerName).Error; err != nil {
		return nil, err
	}
	if provider.Type == "sms" {
		var coon model.SmsConnector
		if err := global.DB.First(&coon, "tenant_id = ?", tenantId).Error; err != nil {
			return nil, err
		}
		return &ProviderResetSms{Config: model.ProviderSms{
			SmsConnectorId: coon.Id,
			SmsConnector:   coon,
			TenantId:       provider.TenantId,
		}}, nil
	}

	return nil, errors.New("provider config not found")
}

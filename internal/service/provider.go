package service

import (
	"accounts/internal/model"
)

func DeleteProvider(id uint) error {
	delList := []any{
		model.ProviderDingTalk{},
		model.ProviderWeCom{},
		model.ProviderOAuth2{},
	}
	if err := deleteSource(model.Provider{}, delList, id, "provider_id"); err != nil {
		return err
	}

	return nil
}

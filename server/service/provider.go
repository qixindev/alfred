package service

import "accounts/models"

func DeleteProvider(id uint) error {
	delList := []any{
		models.ProviderDingTalk{},
		models.ProviderWeCom{},
		models.ProviderOAuth2{},
	}
	if err := deleteSource(models.Provider{}, delList, id, "provider_id"); err != nil {
		return err
	}

	return nil
}

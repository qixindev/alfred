package auth

import (
	"accounts/data"
	"accounts/models"
)

func GetOAuth2User(provider models.Provider) (*UserInfo, error) {
	var config models.ProviderOAuth2Config
	if err := data.DB.First(&config, "tenant_id = ? AND providerId = ?", provider.TenantId, provider.Id).Error; err != nil {
		return nil, err
	}
	return &UserInfo{Name: "oauth2"}, nil
}

package auth

import (
	"accounts/internal/global"
	"accounts/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
)

type Provider interface {
	// Auth Get to external auth. Return redirect location.
	Auth(string) (string, error)

	// Login Callback when auth completed.
	Login(*gin.Context) (*models.UserInfo, error)

	LoginConfig() *gin.H

	ProviderConfig() *gin.H
}

func GetAuthProvider(tenantId uint, providerName string) (Provider, error) {
	var provider models.Provider
	if err := global.DB.First(&provider, "tenant_id = ? AND name = ?", tenantId, providerName).Error; err != nil {
		return nil, err
	}
	if provider.Type == "oauth2" {
		var config models.ProviderOAuth2
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderOAuth2{Config: config}
		return p, nil
	}
	if provider.Type == "dingtalk" {
		var config models.ProviderDingTalk
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderDingTalk{Config: config}
		return p, nil
	}
	if provider.Type == "wecom" {
		var config models.ProviderWeCom
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderWeCom{Config: config}
		return p, nil
	}
	return nil, errors.New("provider config not found")
}

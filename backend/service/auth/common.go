package auth

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
	"github.com/gin-gonic/gin"
)

type ProviderLogin struct {
	State     string `json:"state"`
	AuthState string `json:"authState"`
	Type      string `json:"type"`
	Provider  string `json:"provider"`
	Redirect  string `json:"redirect"`
	ClientId  string `json:"clientId"`
	Tenant    string `json:"tenant"`
	TenantId  uint   `json:"tenantId"`
}

type Provider interface {
	// Auth Get to external auth. Return redirect location.
	Auth(string, string, uint) (string, error)

	// Login Callback when auth completed.
	Login(string, ProviderLogin) (*model.UserInfo, error)

	LoginConfig() *gin.H

	ProviderConfig() *gin.H
}

func GetAuthProvider(tenantId uint, providerName string) (*model.Provider, Provider, error) {
	var provider model.Provider
	if err := global.DB.First(&provider, "tenant_id = ? AND name = ?", tenantId, providerName).Error; err != nil {
		return nil, nil, err
	}
	if provider.Type == "oauth2" {
		var config model.ProviderOAuth2
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderOAuth2{Config: config}
		return &provider, p, nil
	}
	if provider.Type == "dingtalk" {
		var config model.ProviderDingTalk
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderDingTalk{Config: config}
		return &provider, p, nil
	}
	if provider.Type == "wecom" {
		var config model.ProviderWeCom
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderWeCom{Config: config}
		return &provider, p, nil
	}
	if provider.Type == "wechat" {
		var config model.ProviderWechat
		if err := global.DB.First(&config, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, nil, err
		}
		config.Provider.Name = provider.Name
		config.Provider.Type = provider.Type
		p := ProviderWechat{Config: config}
		return &provider, p, nil
	}
	if provider.Type == "sms" {
		var coon model.SmsConnector
		if err := global.DB.First(&coon, "tenant_id = ?", tenantId).Error; err != nil {
			return nil, nil, err
		}
		return &provider, &ProviderSms{Config: model.ProviderSms{
			SmsConnectorId: coon.Id,
			SmsConnector:   coon,
			TenantId:       provider.TenantId,
		}}, nil
	}

	return nil, nil, errors.New("provider config not found")
}

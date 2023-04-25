package models

import "accounts/models/dto"

type ProviderAzureAd struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	TenantId uint `json:"tenantId" gorm:"primaryKey"`
}

type ProviderDingTalk struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	AgentId   string `json:"agentId"`
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`

	TenantId uint `json:"tenantId" gorm:"primaryKey"`
}

func (p *ProviderDingTalk) Dto() dto.ProviderConfigDto {
	return dto.ProviderConfigDto{
		ProviderId:   p.ProviderId,
		Name:         p.Provider.Name,
		Type:         p.Provider.Type,
		AgentId:      p.AgentId,
		ClientId:     p.AppKey,
		ClientSecret: p.AppSecret,
	}
}

type ProviderOAuth2 struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	TokenEndpoint     string `json:"tokenEndpoint"`
	UserinfoEndpoint  string `json:"userinfoEndpoint"`
	Scope             string `json:"scope"`
	ResponseType      string `json:"responseType"`

	TenantId uint `json:"tenantId" gorm:"primaryKey"`
}

func (p *ProviderOAuth2) Dto() dto.ProviderConfigDto {
	return dto.ProviderConfigDto{
		ProviderId:   p.ProviderId,
		Name:         p.Provider.Name,
		Type:         p.Provider.Type,
		ClientId:     p.ClientId,
		ClientSecret: p.ClientSecret,
	}
}

type ProviderWeCom struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	CorpId    string `json:"corpId"`
	AgentId   string `json:"agentId"`
	AppSecret string `json:"appSecret"`

	TenantId uint `json:"tenantId" gorm:"primaryKey"`
}

func (p *ProviderWeCom) Dto() dto.ProviderConfigDto {
	return dto.ProviderConfigDto{
		ProviderId:   p.ProviderId,
		Name:         p.Provider.Name,
		Type:         p.Provider.Type,
		AgentId:      p.CorpId,
		ClientId:     p.AgentId,
		ClientSecret: p.AppSecret,
	}
}

type ProviderSms struct {
	Id             uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId     uint         `json:"providerId"`
	Provider       Provider     `gorm:"foreignKey:ProviderId, TenantId"`
	SmsConnectorId uint         `json:"smsConnectorId"`
	SmsConnector   SmsConnector `gorm:"foreignKey:SmsConnectorId, TenantId"`
	TenantId       uint         `gorm:"primaryKey;autoIncrement" json:"tenantId"`
}

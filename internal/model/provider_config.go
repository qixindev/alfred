package model

import (
	"alfred/internal/endpoint/req"
	"github.com/gin-gonic/gin"
)

type ItfProvider interface {
	Dto() any
	Save(req.Provider) any
}

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

func (p *ProviderDingTalk) Dto() any {
	return &gin.H{
		"providerId": p.ProviderId,
		"name":       p.Provider.Name,
		"type":       p.Provider.Type,
		"agentId":    p.AgentId,
		"appKey":     p.AppKey,
		"appSecret":  p.AppSecret,
	}
}
func (p *ProviderDingTalk) Save(r req.Provider) any {
	return &ProviderDingTalk{
		ProviderId: r.ProviderId,
		TenantId:   r.TenantId,
		AgentId:    r.AgentId,
		AppKey:     r.AppKey,
		AppSecret:  r.AppSecret,
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

func (p *ProviderOAuth2) Dto() any {
	return &gin.H{
		"providerId":        p.ProviderId,
		"name":              p.Provider.Name,
		"type":              p.Provider.Type,
		"clientId":          p.ClientId,
		"scope":             p.Scope,
		"responseType":      p.ResponseType,
		"authorizeEndpoint": p.AuthorizeEndpoint,
		"userInfoEndpoint":  p.UserinfoEndpoint,
		"tokenEndpoint":     p.TokenEndpoint,
	}
}
func (p *ProviderOAuth2) Save(r req.Provider) any {
	return &ProviderOAuth2{
		ProviderId:        r.ProviderId,
		TenantId:          r.TenantId,
		ClientId:          r.ClientId,
		ClientSecret:      r.ClientSecret,
		AuthorizeEndpoint: r.AuthorizeEndpoint,
		UserinfoEndpoint:  r.UserinfoEndpoint,
		TokenEndpoint:     r.TokenEndpoint,
		ResponseType:      r.ResponseType,
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

func (p *ProviderWeCom) Dto() any {
	return &gin.H{
		"providerId": p.ProviderId,
		"name":       p.Provider.Name,
		"type":       p.Provider.Type,
		"corpId":     p.CorpId,
		"agentId":    p.AgentId,
		"appSecret":  p.AppSecret,
	}
}
func (p *ProviderWeCom) Save(r req.Provider) any {
	return &ProviderWeCom{
		ProviderId: r.ProviderId,
		TenantId:   r.TenantId,
		CorpId:     r.CorpId,
		AgentId:    r.AgentId,
		AppSecret:  r.AppSecret,
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

func (p *ProviderSms) Dto() any {
	return &gin.H{
		"providerId": p.ProviderId,
		"name":       p.Provider.Name,
		"type":       p.Provider.Type,
	}
}
func (p *ProviderSms) Save(r req.Provider) any {
	return &ProviderSms{
		ProviderId:     r.ProviderId,
		SmsConnectorId: r.SmsConnectorId,
		TenantId:       r.TenantId,
	}
}

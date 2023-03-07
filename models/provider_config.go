package models

import (
	"accounts/auth"
	"fmt"
	"net/url"
)

type ProviderDingDing struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderWeCom struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderAzureAd struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderOAuth2 struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	TokenEndpoint     string `json:"tokenEndpoint"`
	UserinfoEndpoint  string `json:"userinfoEndpoint"`
	Scope             string `json:"scope"`
	ResponseType      string `json:"responseType"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

func (ProviderOAuth2) TableName() string {
	return "provider_oauth2"
}

func (p ProviderOAuth2) Auth(redirectUri string) string {
	var query url.Values
	query.Set("client_id", p.ClientId)
	query.Set("scope", p.Scope)
	query.Set("response_type", p.ResponseType)
	query.Set("redirect_uri", redirectUri)
	location := fmt.Sprintf("%s?%s", p.AuthorizeEndpoint, query.Encode())
	return location
}

func (ProviderOAuth2) Login() (*auth.UserInfo, error) {
	return &auth.UserInfo{Name: "oauth2"}, nil
}

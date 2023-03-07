package models

type ProviderDingDingConfig struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderWeComConfig struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderAzureAdConfig struct {
	Id         uint     `gorm:"primaryKey" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `json:"provider"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderOAuth2Config struct {
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

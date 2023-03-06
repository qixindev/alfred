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

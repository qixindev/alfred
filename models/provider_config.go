package models

type ProviderDingDingConfig struct {
	Id         uint `gorm:"primaryKey"`
	ProviderId uint
	Provider   Provider

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderWeComConfig struct {
	Id         uint `gorm:"primaryKey"`
	ProviderId uint
	Provider   Provider

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderAzureAdConfig struct {
	Id         uint `gorm:"primaryKey"`
	ProviderId uint
	Provider   Provider

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

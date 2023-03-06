package models

type ProviderDingDingConfig struct {
	Id         uint
	ProviderId uint
	Provider   Provider

	TenantId uint
	Tenant   Tenant
}

type ProviderWeComConfig struct {
	Id         uint
	ProviderId uint
	Provider   Provider

	TenantId uint
	Tenant   Tenant
}

type ProviderAzureAdConfig struct {
	Id         uint
	ProviderId uint
	Provider   Provider

	TenantId uint
	Tenant   Tenant
}

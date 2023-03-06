package models

type Provider struct {
	Id   uint
	Name string
	Type string

	TenantId uint
	Tenant   Tenant
}

type ProviderUser struct {
	Id         uint
	ProviderId uint
	Provider   Provider
	UserId     uint
	User       User
	Name       string

	TenantId uint
	Tenant   Tenant
}

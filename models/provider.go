package models

type Provider struct {
	Id   uint `gorm:"primaryKey"`
	Name string
	Type string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type ProviderUser struct {
	Id         uint `gorm:"primaryKey"`
	ProviderId uint
	Provider   Provider
	UserId     uint
	User       User
	Name       string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

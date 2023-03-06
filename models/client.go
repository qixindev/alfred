package models

type Client struct {
	Id       uint `gorm:"primaryKey"`
	Name     string
	ClientId string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type RedirectUri struct {
	Id          uint `gorm:"primaryKey"`
	ClientId    uint
	Client      Client
	RedirectUri string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

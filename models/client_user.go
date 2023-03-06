package models

type ClientUser struct {
	Id       uint `gorm:"primaryKey"`
	ClientId uint
	Client   Client
	UserId   uint
	User     User
	Sub      string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

package models

type ClientUser struct {
	Id       uint
	ClientId uint
	Client   Client
	UserId   uint
	User     User
	Sub      string

	TenantId uint
	Tenant   Tenant
}

package models

type Client struct {
	Id       uint
	Name     string
	ClientId string

	TenantId uint
	Tenant   Tenant
}

type RedirectUri struct {
	Id          uint
	ClientId    uint
	Client      Client
	RedirectUri string

	TenantId uint
	Tenant   Tenant
}

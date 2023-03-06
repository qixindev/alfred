package models

type Device struct {
	Id   uint
	Name string

	TenantId uint
	Tenant   Tenant
}

package models

type Group struct {
	Id       uint
	Name     string
	ParentId uint

	TenantId uint
	Tenant   Tenant
}

type GroupUser struct {
	Id      uint
	GroupId uint
	Group   Group
	UserId  uint
	User    User

	TenantId uint
	Tenant   Tenant
}

type GroupDevice struct {
	Id       uint
	GroupId  uint
	Group    Group
	DeviceId uint
	Device   Device

	TenantId uint
	Tenant   Tenant
}

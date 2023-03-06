package models

import "accounts/models/dto"

type Group struct {
	Id       uint `gorm:"primaryKey"`
	Name     string
	ParentId uint

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type GroupUser struct {
	Id      uint `gorm:"primaryKey"`
	GroupId uint
	Group   Group
	UserId  uint
	User    User
	Role    string

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type GroupDevice struct {
	Id       uint `gorm:"primaryKey"`
	GroupId  uint
	Group    Group
	DeviceId uint
	Device   Device

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

func (g *Group) Dto() dto.GroupDto {
	return dto.GroupDto{
		Id:       g.Id,
		Name:     g.Name,
		ParentId: g.ParentId,
	}
}

func Group2Dto(g Group) dto.GroupDto {
	return g.Dto()
}

func (g *Group) GroupMemberDto() dto.GroupMemberDto {
	return dto.GroupMemberDto{
		Type: "group",
		Id:   g.Id,
		Name: g.Name,
	}
}

func (u *GroupUser) GroupMemberDto() dto.GroupMemberDto {
	return dto.GroupMemberDto{
		Type: "user",
		Id:   u.Id,
		Name: u.User.Username,
		Role: u.Role,
	}
}

func (d *GroupDevice) GroupMemberDto() dto.GroupMemberDto {
	return dto.GroupMemberDto{
		Type: "device",
		Id:   d.Id,
		Name: d.Device.Name,
	}
}

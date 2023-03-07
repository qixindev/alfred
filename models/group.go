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
	Id      uint   `gorm:"primaryKey" json:"id"`
	GroupId uint   `json:"groupId"`
	Group   Group  `gorm:"foreignKey:GroupId, TenantId" json:"group"`
	UserId  uint   `json:"userId"`
	User    User   `gorm:"foreignKey:UserId, TenantId" json:"user"`
	Role    string `json:"role"`

	TenantId uint `gorm:"primaryKey"`
}

type GroupDevice struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupId  uint   `json:"groupId"`
	Group    Group  `gorm:"foreignKey:GroupId, TenantId" json:"group"`
	DeviceId uint   `json:"deviceId"`
	Device   Device `gorm:"foreignKey:DeviceId, TenantId" json:"device"`

	TenantId uint `gorm:"primaryKey"`
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

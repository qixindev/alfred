package models

import (
	"accounts/models/dto"
)

type ResourceType struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ClientId uint   `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	TenantId uint   `gorm:"primaryKey"`
}

type Resource struct {
	Id       uint         `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   uint         `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"type"`
	ParentId uint         `json:"parent"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeAction struct {
	Id       uint         `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   uint         `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"type"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeRole struct {
	Id       uint         `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   uint         `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"type"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeRoleAction struct {
	Id         uint               `gorm:"primaryKey" json:"id"`
	RoleId     uint               `json:"roleId"`
	Role       ResourceTypeRole   `gorm:"foreignKey:RoleId, TenantId" json:"role"`
	ActionId   uint               `json:"actionId"`
	Action     ResourceTypeAction `gorm:"foreignKey:ActionId, TenantId" json:"action"`
	ActionName string             `gorm:"column:action_name" json:"actionName"`
	TenantId   uint               `gorm:"primaryKey"`
}

func (r *ResourceTypeRoleAction) Dto() *dto.ResourceTypeRoleActionDto {
	return &dto.ResourceTypeRoleActionDto{
		Id:         r.Id,
		RoleId:     r.RoleId,
		TenantId:   r.TenantId,
		ActionId:   r.ActionId,
		ActionName: r.ActionName,
	}
}

type ResourceRoleUser struct {
	Id           uint             `gorm:"primaryKey" json:"id"`
	ResourceId   uint             `json:"resourceId"`
	Resource     Resource         `gorm:"foreignKey:ResourceId, TenantId" json:"resource"`
	RoleId       uint             `json:"roleId"`
	Role         ResourceTypeRole `gorm:"foreignKey:RoleId, TenantId" json:"role"`
	ClientUserId uint             `json:"userId"`
	ClientUser   ClientUser       `gorm:"foreignKey:ClientUserId, TenantId" json:"user"`
	TenantId     uint             `gorm:"primaryKey"`
}

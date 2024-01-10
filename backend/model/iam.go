package model

import (
	"alfred/backend/endpoint/dto"
)

type ResourceType struct {
	Id       string `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"-"`
	TenantId uint   `gorm:"primaryKey"`
}

type Resource struct {
	Id       string       `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   string       `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"-"`
	ParentId string       `json:"parent"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeAction struct {
	Id       string       `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   string       `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"-"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeRole struct {
	Id       string       `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	TypeId   string       `json:"typeId"`
	Type     ResourceType `gorm:"foreignKey:TypeId, TenantId" json:"-"`
	TenantId uint         `gorm:"primaryKey"`
}

type ResourceTypeRoleAction struct {
	Id         uint               `gorm:"primaryKey" json:"id"`
	RoleId     string             `json:"roleId"`
	Role       ResourceTypeRole   `gorm:"foreignKey:RoleId, TenantId" json:"-"`
	ActionId   string             `json:"actionId"`
	Action     ResourceTypeAction `gorm:"foreignKey:ActionId, TenantId" json:"-"`
	ActionName string             `gorm:"<-:false;-:migration" json:"actionName"`
	RoleName   string             `gorm:"<-:false;-:migration" json:"roleName"`
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
func ResourceRoleActionDto(r ResourceTypeRoleAction) *dto.ResourceTypeRoleActionDto {
	return r.Dto()
}

type ResourceRoleUser struct {
	Id           uint             `gorm:"primaryKey" json:"id"`
	ResourceId   string           `json:"resourceId"`
	Resource     Resource         `gorm:"foreignKey:ResourceId, TenantId" json:"resource"`
	ResourceName string           `json:"resourceName" gorm:"<-:false;-:migration"`
	RoleId       string           `json:"roleId"`
	Role         ResourceTypeRole `gorm:"foreignKey:RoleId, TenantId" json:"role"`
	RoleName     string           `json:"roleName" gorm:"<-:false;-:migration"`
	ClientUserId uint             `json:"userId"`
	ClientUser   ClientUser       `gorm:"foreignKey:ClientUserId, TenantId" json:"user"`
	TenantId     uint             `gorm:"primaryKey"`
	Sub          string           `json:"sub" gorm:"<-:false;-:migration"`
	DisplayName  string           `json:"displayName" gorm:"<-:false;-:migration"`
}

func (r *ResourceRoleUser) Dto() *dto.ResourceRoleUserDto {
	return &dto.ResourceRoleUserDto{
		Id:           r.Id,
		ResourceId:   r.ResourceId,
		ResourceName: r.ResourceName,
		RoleId:       r.RoleId,
		RoleName:     r.RoleName,
		Sub:          r.Sub,
		DisplayName:  r.DisplayName,
		ClientUserId: r.ClientUserId,
	}
}

func ResourceRoleUserDto(r ResourceRoleUser) *dto.ResourceRoleUserDto {
	return r.Dto()
}

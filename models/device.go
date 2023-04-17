package models

import (
	"accounts/models/dto"
	"time"
)

type Device struct {
	Id   string `gorm:"primaryKey;not null" json:"id"`
	Name string `json:"name"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

func Device2Dto(d Device) dto.DeviceDto {
	return d.Dto()
}

func (d *Device) Dto() dto.DeviceDto {
	return dto.DeviceDto{
		Id:   d.Id,
		Name: d.Name,
	}
}

type DeviceSecret struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	DeviceId string `json:"deviceId"`
	Device   Device `gorm:"foreignKey:DeviceId, TenantId" json:"device"`
	Secret   string `json:"secret"`
	TenantId uint   `gorm:"primaryKey" json:"tenantId"`
}

func (d *DeviceSecret) Dto() dto.DeviceSecretDto {
	return dto.DeviceSecretDto{
		Id:     d.Id,
		Name:   d.Name,
		Secret: d.Secret,
	}
}

func DeviceSecret2Dto(s DeviceSecret) dto.DeviceSecretDto {
	return s.Dto()
}

type DeviceCode struct {
	Id        uint      `json:"id"`
	TenantId  uint      `json:"tenantId"`
	Tenant    Tenant    `gorm:"foreignKey:TenantId" json:"tenant"`
	Code      string    `json:"code" gorm:"uniqueIndex"`
	UserCode  string    `json:"userCode" gorm:"uniqueIndex"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

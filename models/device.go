package models

import "accounts/models/dto"

type Device struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
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

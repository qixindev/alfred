package models

import "accounts/models/dto"

type Tenant struct {
	Id   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

func (t *Tenant) Dto() dto.TenantDto {
	return dto.TenantDto{
		Id:   t.Id,
		Name: t.Name,
	}
}

func Tenant2Dto(t Tenant) dto.TenantDto {
	return t.Dto()
}

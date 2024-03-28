package model

import (
	"alfred/backend/endpoint/dto"
)

type Tenant struct {
	Id        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `json:"name" gorm:"uniqueIndex"`
	LoginPage string `json:"loginPage,omitempty" gorm:"type:jsonb;default:'{}'"`
	Proto     string `json:"proto,omitempty" gorm:"type:jsonb;default:'[]'"`
	Sub       string `json:"sub,omitempty" gorm:"<-:false;-:migration"`
	Role      string `json:"role,omitempty" gorm:"<-:false;-:migration"`
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

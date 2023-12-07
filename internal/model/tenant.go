package model

import (
	"alfred/internal/endpoint/dto"
)

type Tenant struct {
	Id        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `json:"name" gorm:"uniqueIndex"`
	LoginPage string `json:"loginPage" gorm:"type:jsonb;default:'{}'"`
	Proto     string `json:"proto" gorm:"type:jsonb;default:'[]'"`
	Sub       string `json:"sub" gorm:"<-:false;-:migration"`
	Role      string `json:"role" gorm:"<-:false;-:migration"`
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

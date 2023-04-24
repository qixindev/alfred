package models

import (
	"accounts/models/dto"
)

type Provider struct {
	Id           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type" gorm:"not null"`
	AgentId      string `json:"agentId" gorm:"<-:false;-:migration"`
	ClientId     string `json:"clientId" gorm:"<-:false;-:migration"`
	ClientSecret string `json:"clientSecret" gorm:"<-:false;-:migration"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

func (p *Provider) Dto() dto.ProviderDto {
	return dto.ProviderDto{
		Id:   p.Id,
		Name: p.Name,
		Type: p.Type,
	}
}

func Provider2Dto(p Provider) dto.ProviderDto {
	return p.Dto()
}

type ProviderUser struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`
	UserId     uint     `json:"userId"`
	User       User     `gorm:"foreignKey:UserId, TenantId" json:"user"`
	Name       string   `json:"name"`

	TenantId uint `gorm:"primaryKey"`
}

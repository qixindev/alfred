package models

import (
	"accounts/models/dto"
)

type Client struct {
	Id   string `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type RedirectUri struct {
	Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientId    string `json:"clientId"`
	Client      Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	RedirectUri string `json:"redirectUri"`

	TenantId uint `gorm:"primaryKey"`
}

type ClientSecret struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	Secret   string `json:"secret"`

	TenantId uint `gorm:"primaryKey"`
}

func (c *Client) Dto() dto.ClientDto {
	return dto.ClientDto{
		Id:   c.Id,
		Name: c.Name,
	}
}

func Client2Dto(c Client) dto.ClientDto {
	return c.Dto()
}

func (r *RedirectUri) Dto() dto.RedirectUriDto {
	return dto.RedirectUriDto{
		Id:          r.Id,
		RedirectUri: r.RedirectUri,
	}
}

func RedirectUri2Dto(r RedirectUri) dto.RedirectUriDto {
	return r.Dto()
}

func (s *ClientSecret) Dto() dto.ClientSecretDto {
	return dto.ClientSecretDto{
		Id:     s.Id,
		Name:   s.Name,
		Secret: s.Secret,
	}
}

func ClientSecret2Dto(s ClientSecret) dto.ClientSecretDto {
	return s.Dto()
}

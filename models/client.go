package models

import "accounts/models/dto"

type Client struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

type RedirectUri struct {
	Id          uint   `gorm:"primaryKey" json:"id"`
	ClientId    uint   `json:"clientId"`
	Client      Client `json:"client"`
	RedirectUri string `json:"redirectUri"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

func (c *Client) Dto() dto.ClientDto {
	return dto.ClientDto{
		Id:       c.Id,
		Name:     c.Name,
		ClientId: c.ClientId,
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

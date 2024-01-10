package model

import (
	"alfred/backend/endpoint/dto"
)

type ClientUser struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientId string `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	UserId   uint   `json:"userId"`
	User     User   `gorm:"foreignKey:UserId, TenantId" json:"user"`
	Sub      string `json:"sub"`
	UserName string `json:"userName" gorm:"<-:false;-:migration"`
	Phone    string `json:"phone" gorm:"<-:false;-:migration"`
	Email    string `json:"email" gorm:"<-:false;-:migration"`

	TenantId uint `gorm:"primaryKey"`
}

func (c *ClientUser) Dto() dto.ClientUserDto {
	return dto.ClientUserDto{
		Id:       c.Id,
		Sub:      c.Sub,
		ClientId: c.ClientId,
		UserName: c.UserName,
		Phone:    c.Phone,
		Email:    c.Email,
	}
}

func ClientUserDto(c ClientUser) dto.ClientUserDto {
	return c.Dto()
}

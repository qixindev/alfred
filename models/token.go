package models

import (
	"time"
)

type TokenCode struct {
	Id        uint      `json:"id"`
	Token     string    `json:"token"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`

	ClientId uint   `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	TenantId uint   `json:"tenantId"`
}

package model

import (
	"time"
)

type TokenCode struct {
	Id        uint      `json:"id"`
	Token     string    `json:"token"`
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	Sub       string    `json:"sub"`
	CreatedAt time.Time `json:"createdAt"`

	ClientId string `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	TenantId uint   `json:"tenantId"`
}

type Paging struct {
	Search   string `json:"-" gorm:"-" form:"search"`
	PageNum  int    `json:"-" gorm:"-" form:"pageNum"`
	PageSize int    `json:"-" gorm:"-" form:"pageSize"`
}

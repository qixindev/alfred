package models

type ClientUser struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	ClientId uint   `json:"clientId"`
	Client   Client `json:"client"`
	UserId   uint   `json:"userId"`
	User     User   `json:"user"`
	Sub      string `json:"sub"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
}

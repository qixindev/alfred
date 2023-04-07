package models

type ClientUser struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientId string `json:"clientId"`
	Client   Client `gorm:"foreignKey:ClientId, TenantId" json:"client"`
	UserId   uint   `json:"userId"`
	User     User   `gorm:"foreignKey:UserId, TenantId" json:"user"`
	Sub      string `json:"sub"`

	TenantId uint `gorm:"primaryKey"`
}

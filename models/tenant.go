package models

type Tenant struct {
	Id   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

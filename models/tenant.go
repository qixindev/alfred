package models

type Tenant struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

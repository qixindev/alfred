package models

type Tenant struct {
	Id   uint `gorm:"primaryKey"`
	Name string
}

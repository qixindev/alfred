package data

import (
	"accounts/auth"
	"accounts/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func migrateDB() {
	DB.AutoMigrate(&models.Tenant{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Group{})
	DB.AutoMigrate(&models.Client{})
	DB.AutoMigrate(&models.ClientUser{})
	DB.AutoMigrate(&models.Device{})
	DB.AutoMigrate(&models.GroupUser{})
	DB.AutoMigrate(&models.GroupDevice{})
	DB.AutoMigrate(&models.RedirectUri{})
	DB.AutoMigrate(&models.ClientSecret{})
	DB.AutoMigrate(&models.TokenCode{})
	DB.AutoMigrate(&models.Provider{})
	DB.AutoMigrate(&models.ProviderUser{})

	DB.AutoMigrate(&auth.ProviderOAuth2{})
	DB.AutoMigrate(&auth.ProviderDingTalk{})
}

func InitDB() error {
	dsn := os.Getenv("dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	migrateDB()
	return nil
}

package data

import (
	"accounts/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func migrateDB() {
	//DB.AutoMigrate(&models.Tenant{})
	//DB.AutoMigrate(&models.User{})
	//DB.AutoMigrate(&models.Group{})
	//DB.AutoMigrate(&models.Client{})
	//DB.AutoMigrate(&models.ClientUser{})
	//DB.AutoMigrate(&models.Device{})
	//DB.AutoMigrate(&models.GroupUser{})
	//DB.AutoMigrate(&models.GroupDevice{})
	//DB.AutoMigrate(&models.RedirectUri{})
	//DB.AutoMigrate(&models.ClientSecret{})
	//DB.AutoMigrate(&models.TokenCode{})
	//DB.AutoMigrate(&models.ProviderUser{})
	//
	//DB.AutoMigrate(&models.Provider{})
	//DB.AutoMigrate(&models.ProviderOAuth2{})
	//DB.AutoMigrate(&models.ProviderDingTalk{})
	//DB.AutoMigrate(&models.ProviderWeCom{})
	//
	//DB.AutoMigrate(&models.SmsConnector{})
	//DB.AutoMigrate(&models.SmsTcloud{})
	//
	//DB.AutoMigrate(&models.ResourceType{})
	//DB.AutoMigrate(&models.ResourceTypeAction{})
	//DB.AutoMigrate(&models.ResourceTypeRole{})
	//DB.AutoMigrate(&models.ResourceTypeRoleAction{})
	//DB.AutoMigrate(&models.Resource{})
	//DB.AutoMigrate(&models.ResourceRoleUser{})

	//DB.Migrator().CreateTable(&models.ResourceRoleUser{})
}

func CheckFirstRun() error {
	var tenant models.Tenant
	if err := DB.First(&tenant, "name = ?", "default").Error; err != nil {
		tenant.Name = "default"
		if err := DB.Create(&tenant).Error; err != nil {
			return err
		}
	}
	return nil
}

func InitDB() error {
	dsn := os.Getenv("dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	migrateDB()
	err = CheckFirstRun()
	return err
}

func WithTenant(tenantId uint) *gorm.DB {
	return DB.Where("tenant_id = ?", tenantId)
}

package initial

import (
	"accounts/global"
	"accounts/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateDB() {
	migrateList := []any{
		&models.Tenant{},
		&models.User{},
		&models.Group{},
		&models.Client{},
		&models.ClientUser{},
		&models.Device{},
		&models.GroupUser{},
		&models.GroupDevice{},
		&models.RedirectUri{},
		&models.ClientSecret{},
		&models.TokenCode{},
		&models.ProviderUser{},
		&models.Provider{},
		&models.ProviderOAuth2{},
		&models.ProviderDingTalk{},
		&models.ProviderWeCom{},
		&models.SmsConnector{},
		&models.SmsTcloud{},
		&models.ResourceType{},
		&models.ResourceTypeAction{},
		&models.ResourceTypeRole{},
		&models.ResourceTypeRoleAction{},
		&models.Resource{},
		&models.ResourceRoleUser{},
	}

	if err := global.DB.AutoMigrate(migrateList...); err != nil {
		fmt.Println("migrate db err: ", err)
		return
	}
}

func CheckFirstRun() error {
	var tenant models.Tenant
	if err := global.DB.First(&tenant, "name = ?", "default").Error; err != nil {
		tenant.Name = "default"
		if err = global.DB.Create(&tenant).Error; err != nil {
			return err
		}
	}
	return nil
}

func InitDB() error {
	dsn := "host=base.postgres.database.chinacloudapi.cn port=5432 dbname=accounts user=qixin password=Qx20190822"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	global.DB = db
	// migrateDB()
	err = CheckFirstRun()
	return err
}

func WithTenant(tenantId uint) *gorm.DB {
	return global.DB.Where("tenant_id = ?", tenantId)
}

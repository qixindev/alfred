package initial

import (
	"alfred/internal/model"
	"alfred/pkg/global"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() error {
	dsn := global.CONFIG.Pgsql.ConfigDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	global.DB = db
	return nil
}

func InitDefaultTenant() error {
	if global.DB == nil {
		return errors.New("global db is nil")
	}
	tenant := model.Tenant{
		Name: "default",
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Tenant{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&tenant).Error; err != nil {
				return err
			}
		}

		client := model.Client{
			Id:       "default",
			Name:     "default",
			TenantId: tenant.Id,
		}
		if err := tx.First(&model.Client{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&client).Error; err != nil {
				return err
			}
		}

		clientSecret := model.ClientSecret{
			Name:     "default",
			Secret:   "multi-tenant",
			ClientId: client.Id,
			TenantId: tenant.Id,
		}
		if err := tx.First(&model.ClientSecret{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&clientSecret).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

package initial

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"github.com/pkg/errors"
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

const (
	DefaultTenant = "default"
	DefaultClient = "default"
	DefaultUser   = "admin"
	DefaultPwd    = "admin"
)

func InitDefaultTenant() error {
	if global.DB == nil {
		return errors.New("global db is nil")
	}

	tenant := model.Tenant{Name: DefaultTenant}
	client := model.Client{
		Id:   DefaultClient,
		Name: DefaultClient,
	}
	var tmpUser model.User
	return global.DB.Debug().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name = ?", DefaultTenant).First(&tenant).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&tenant).Error; err != nil {
				return errors.New("create tenant err")
			}
		}

		if err := tx.First(&client).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			client.TenantId = tenant.Id
			if err = tx.Create(&client).Error; err != nil {
				return errors.New("create client err")
			}
		}

		if err := tx.First(&model.ClientSecret{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model.ClientSecret{
				Name:     DefaultClient,
				Secret:   "multi-tenant",
				ClientId: client.Id,
				TenantId: tenant.Id,
			}).Error; err != nil {
				return errors.New("create secret err")
			}
		}
		if err := tx.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			adminPwd, err := utils.HashPassword("admin")
			if err != nil {
				return err
			}
			user := model.User{
				Username:         DefaultUser,
				PasswordHash:     adminPwd,
				EmailVerified:    false,
				PhoneVerified:    false,
				TwoFactorEnabled: false,
				Disabled:         false,
				TenantId:         tenant.Id,
				Role:             DefaultPwd,
				Meta:             "{}",
				From:             "init",
				Avatar:           "",
			}
			if err = tx.Create(&user).Error; err != nil {
				return err
			}
			tmpUser.Id = user.Id
		}
		if err := tx.First(&model.ClientUser{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model.ClientUser{
				TenantId: tenant.Id,
				ClientId: client.Id,
				UserId:   tmpUser.Id,
				Sub:      DefaultUser,
			}).Error; err != nil {
				return errors.New("create secret err")
			}
		}

		if _, err := utils.LoadRsaPublicKeys(tenant.Name); err != nil {
			return errors.WithMessage(err, "LoadRsaPublicKeys err")
		}
		return nil
	})
}

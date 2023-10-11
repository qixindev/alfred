package initial

import (
	"alfred/internal/model"
	"alfred/pkg/global"
	"alfred/pkg/utils"
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

	var tmpTenant model.Tenant
	var tmpClient model.Client
	var tmpUser model.User
	//return global.DB.Debug().Transaction(func(_ *gorm.DB) error {
	//	tx := global.DB
	//	if err := tx.First(&tmpTenant).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		tenant := model.Tenant{
	//			Name: "default",
	//		}
	//		if err = tx.Create(&tenant).Error; err != nil {
	//			return errors.New("create tenant err")
	//		}
	//		tmpTenant.Id = tenant.Id
	//	}
	//
	//	if err := tx.First(&tmpClient).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		client := model.Client{
	//			Id:       "default",
	//			Name:     "default",
	//			TenantId: tmpTenant.Id,
	//		}
	//		if err = tx.Create(&client).Error; err != nil {
	//			return errors.New("create client err")
	//		}
	//		tmpClient.Id = client.Id
	//	}
	//
	//	if err := tx.First(&model.ClientSecret{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		if err = tx.Create(&model.ClientSecret{
	//			Name:     "default",
	//			Secret:   "multi-tenant",
	//			ClientId: tmpClient.Id,
	//			TenantId: tmpTenant.Id,
	//		}).Error; err != nil {
	//			return errors.New("create secret err")
	//		}
	//	}
	//	if err := tx.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		adminPwd, err := utils.HashPassword("admin")
	//		if err != nil {
	//			return err
	//		}
	//		user := model.User{
	//			Username:         "admin",
	//			PasswordHash:     adminPwd,
	//			EmailVerified:    false,
	//			PhoneVerified:    false,
	//			TwoFactorEnabled: false,
	//			Disabled:         false,
	//			TenantId:         tmpTenant.Id,
	//			Role:             "admin",
	//			Meta:             "{}",
	//			From:             "create",
	//			Avatar:           "",
	//		}
	//		if err = global.DB.Create(&user).Error; err != nil {
	//			return err
	//		}
	//		tmpUser.Id = user.Id
	//	}
	//	if err := tx.First(&model.ClientUser{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		if err = tx.Create(&model.ClientUser{
	//			TenantId: tmpTenant.Id,
	//			ClientId: tmpClient.Id,
	//			UserId:   tmpUser.Id,
	//			Sub:      "admin",
	//		}).Error; err != nil {
	//			return errors.New("create secret err")
	//		}
	//	}
	//
	//	if _, err := utils.LoadRsaPublicKeys(tmpTenant.Name); err != nil {
	//		return errors.New("LoadRsaPublicKeys err")
	//	}
	//	return nil
	//})
	var err error
	return global.DB.Debug().Transaction(func(tx *gorm.DB) error {
		tenant := model.Tenant{
			Name: "default",
		}
		if err = tx.Create(&tenant).Error; err != nil {
			return errors.New("create tenant err")
		}

		client := model.Client{
			Id:       "default",
			Name:     "default",
			TenantId: tenant.Id,
		}
		if err = tx.Create(&client).Error; err != nil {
			return errors.New("create client err")
		}
		tmpClient.Id = client.Id

		if err = tx.Create(&model.ClientSecret{
			Name:     "default",
			Secret:   "multi-tenant",
			ClientId: client.Id,
			TenantId: tenant.Id,
		}).Error; err != nil {
			return errors.New("create secret err")
		}

		adminPwd, err := utils.HashPassword("admin")
		if err != nil {
			return err
		}
		user := model.User{
			Username:         "admin",
			PasswordHash:     adminPwd,
			EmailVerified:    false,
			PhoneVerified:    false,
			TwoFactorEnabled: false,
			Disabled:         false,
			TenantId:         tenant.Id,
			Role:             "admin",
			Meta:             "{}",
			From:             "create",
			Avatar:           "",
		}
		if err = global.DB.Create(&user).Error; err != nil {
			return err
		}
		tmpUser.Id = user.Id

		if err = tx.Create(&model.ClientUser{
			TenantId: tenant.Id,
			ClientId: client.Id,
			UserId:   user.Id,
			Sub:      "admin",
		}).Error; err != nil {
			return errors.New("create secret err")
		}

		if _, err := utils.LoadRsaPublicKeys(tmpTenant.Name); err != nil {
			return errors.New("LoadRsaPublicKeys err")
		}
		return nil
	})
}

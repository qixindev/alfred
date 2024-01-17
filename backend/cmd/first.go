package cmd

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
)

const (
	DefaultTenant = "default"
	DefaultClient = "default"
	DefaultUser   = "admin"
	DefaultPwd    = "admin"
)

func initFirstRun() {
	if err := initSystem(); err != nil {
		fmt.Println("init system err:", err)
		os.Exit(1)
		return
	}

	var tenant model.Tenant
	if err := global.DB.First(&tenant, "name = ?", DefaultTenant).Error; err == nil {
		fmt.Println("Default tenant is already in use")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("get tenant err: ", err.Error())
		os.Exit(1)
	}

	if err := insertDB(); err != nil {
		fmt.Println("insert database error:", err)
		os.Exit(2)
		return
	}

	if _, err := utils.LoadRsaPublicKeys(DefaultTenant); err != nil {
		fmt.Println("load rsa public keys error:", err)
		os.Exit(2)
		return
	}

	fmt.Println("===== Success =====")
}

func insertDB() error {
	var tenant model.Tenant
	tenant.Name = DefaultTenant
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tenant).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Client{Id: DefaultClient, Name: DefaultClient, TenantId: tenant.Id}).Error; err != nil {
			return err
		}

		adminPwd, err := utils.HashPassword(DefaultPwd)
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
			Role:             "admin",
			From:             "init",
			Meta:             "{}",
		}
		if err = tx.Create(&user).Error; err != nil {
			return err
		}

		if err = tx.Create(&model.ClientUser{
			ClientId: DefaultClient,
			TenantId: tenant.Id,
			Sub:      uuid.NewString(),
			UserId:   user.Id,
		}).Error; err != nil {
			return err
		}

		if err = tx.Create(&model.ClientSecret{
			ClientId: DefaultClient,
			Name:     DefaultClient,
			Secret:   uuid.NewString(),
			TenantId: tenant.Id,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

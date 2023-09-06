package cmd

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"errors"
	"fmt"
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
	var tenant model.Tenant
	if err := global.DB.First(&tenant, "name = ?", DefaultTenant).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return
	} else if err != nil {
		return
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
	if err := global.DB.Create(&tenant).Error; err != nil {
		return err
	}

	if err := global.DB.Create(&model.Client{Id: DefaultClient, Name: DefaultClient, TenantId: tenant.Id}).Error; err != nil {
		return err
	}

	adminPwd, err := utils.HashPassword(DefaultPwd)
	if err != nil {
		return err
	}
	if err = global.DB.Create(&model.User{
		Username:         DefaultUser,
		PasswordHash:     adminPwd,
		EmailVerified:    false,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		Disabled:         false,
		TenantId:         tenant.Id,
		Role:             "admin",
	}).Error; err != nil {
		return err
	}

	return nil
}

package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
	"gorm.io/gorm"
)

func DeleteTenant(tenant model.Tenant) error {
	if tenant.Id == 0 {
		return errors.New("delete invalid tenant")
	}
	delList := []any{
		model.RedirectUri{}, model.TokenCode{},
		model.ResourceTypeRoleAction{}, model.ResourceRoleUser{},
		model.ResourceTypeAction{}, model.Resource{}, model.ResourceTypeRole{}, model.ResourceType{},

		model.ProviderUser{}, model.ProviderDingTalk{}, model.ProviderWeCom{},
		model.ProviderSms{}, model.ProviderOAuth2{}, model.Provider{},

		model.GroupUser{}, model.GroupDevice{}, model.Group{},
		model.DeviceSecret{}, model.DeviceCode{}, model.Device{},
		model.ClientUser{}, model.ClientSecret{}, model.Client{},
		model.User{},
	}

	// 开启数据库事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		for _, v := range delList {
			if err := tx.Model(v).Where("tenant_id = ?", tenant.Id).Delete(v).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := tx.Where("id = ?", tenant.Id).Delete(model.Tenant{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

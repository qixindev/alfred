package service

import (
	"accounts/internal/model"
	"accounts/pkg/global"
)

func CopyUser(sub string, tenantId uint) error {
	var clientUser model.ClientUser
	if err := global.DB.Model(clientUser).Where("sub = ?", sub).First(&clientUser).Error; err != nil {
		return err
	}

	sql := `INSERT INTO users (username, first_name, last_name, display_name, email, email_verified,
                   password_hash, phone, phone_verified, two_factor_enabled, disabled, tenant_id, role)
			SELECT username, first_name, last_name, display_name, email, email_verified,
       				password_hash, phone, phone_verified, two_factor_enabled, disabled, ? as tenant_id, 'owner' as role 
			FROM users WHERE id = ?;`
	if err := global.DB.Exec(sql, tenantId, clientUser.UserId).Error; err != nil {
		return err
	}

	return nil
}

func DeleteUser(id uint) error {
	var clientUser []uint
	if err := global.DB.Model(model.ClientUser{}).Select("id").Where("user_id = ?", id).Find(&clientUser).Error; err == nil {
		if err = global.DB.Where("client_user_id in ?", clientUser).
			Delete(model.ResourceRoleUser{}).Error; err != nil {
			return err
		}
	}

	if err := global.DB.Where("user_id = ?", id).Delete(&model.ClientUser{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("user_id = ?", id).Delete(&model.GroupUser{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("user_id = ?", id).Delete(&model.ProviderUser{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}

	return nil
}

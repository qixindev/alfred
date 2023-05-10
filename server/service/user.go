package service

import (
	"accounts/global"
	"accounts/models"
)

func CopyUser(sub string, tenantId uint) error {
	var clientUser models.ClientUser
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

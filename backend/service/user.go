package service

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func CreateUser(user model.User) (*model.User, error) {
	if user.DisplayName == "" || user.Username == "" {
		return nil, errors.New("invalid user parameter")
	}
	if err := global.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByPhone(phone string, tenantId uint) (*model.User, error) {
	var user model.User
	if err := global.DB.Where("tenant_id = ? AND phone = ?", tenantId, phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByEmail(email string, tenantId uint) (*model.User, error) {
	var user model.User
	if err := global.DB.Where("tenant_id AND email = ?", tenantId, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func BindLoginUser(userInfo *model.UserInfo, tenantId uint, userFrom string) (user *model.User, err error) {
	newUser := model.User{
		Username:         userInfo.Name,
		FirstName:        userInfo.FirstName,
		LastName:         userInfo.LastName,
		DisplayName:      userInfo.DisplayName,
		Email:            userInfo.Email,
		EmailVerified:    false,
		Phone:            userInfo.Phone,
		PhoneVerified:    false,
		TwoFactorEnabled: false,
		Disabled:         false,
		TenantId:         tenantId,
		From:             userFrom,
		Avatar:           userInfo.Picture,
		Meta:             "{}",
	}
	if newUser.Username == "" {
		newUser.Username = uuid.NewString()
	}
	if userInfo.Phone == "" && userInfo.Email == "" {
		global.LOG.Info("create user: " + userInfo.Name + " " + userInfo.DisplayName)
		return CreateUser(newUser) // 无需绑定，直接创建
	}

	if userInfo.Email != "" {
		user, err = GetUserByEmail(userInfo.Email, tenantId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.LOG.Info("no such email user, creating")
			return CreateUser(newUser)
		} else if err != nil {
			return nil, err
		}
		global.LOG.Info("bind email user: " + user.Email)
		return user, nil
	}

	user, err = GetUserByPhone(userInfo.Phone, tenantId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Info("no such phone user, creating")
		return CreateUser(newUser)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserBySubId(tenantId uint, clientId string, subId string) (*model.User, error) {
	if tenantId == 0 || clientId == "" || subId == "" {
		return nil, errors.New("invalidate GetUserBySubId param")
	}
	var user model.User
	if err := global.DB.Table("users as u").
		Select("u.id", "u.username", "u.display_name", "u.email", "u.phone", "u.disabled", "u.role", "u.avatar", "u.from",
			"u.password_hash", "u.tenant_id", "u.email_verified", "u.phone_verified").
		Joins("LEFT JOIN client_users as cu ON cu.user_id = u.id").
		Where("cu.sub = ? AND cu.client_id = ? AND u.tenant_id = ?", subId, clientId, tenantId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUserPhoneOrEmailExist(user model.User, tenantId uint) (bool, error) {
	if user.Phone == "" && user.Email == "" {
		return false, nil
	}
	if user.Phone != "" {
		if err := global.DB.Model(user).Where("tenant_id = ? AND phone = ?", tenantId, user.Phone).
			First(&model.User{}).Error; err == nil {
			return true, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}
	if user.Email != "" {
		if err := global.DB.Model(user).Where("tenant_id = ? AND email = ?", tenantId, user.Email).
			First(&model.User{}).Error; err == nil {
			return true, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}
	return false, nil
}

func GetTenantIdByTenantName(tenantName string) (uint, error) {
	var tenant model.Tenant
	if err := global.DB.Model(tenant).Where("name = ?", tenantName).First(&tenant).Error; err != nil {
		return 0, err
	}
	return tenant.Id, nil
}

func GetAlfredClientUser(clientId string, tenantId uint, userId uint) (string, error) {
	var clientUser model.ClientUser
	if err := global.DB.Model(clientUser).Where("client_id = ? AND tenant_id = ? AND user_id = ?", clientId, tenantId, userId).
		First(&clientUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		clientUser.UserId = userId
		clientUser.ClientId = clientId
		clientUser.TenantId = tenantId
		clientUser.Sub = uuid.NewString()
		if err = global.DB.Create(&clientUser).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return clientUser.Sub, nil
}

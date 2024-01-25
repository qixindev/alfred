package rg

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"gorm.io/gorm"
)

func GetResourceGroupRoleList(tenantId uint, groupId string) ([]model.ResourceGroupRole, error) {
	var roles []model.ResourceGroupRole
	if err := global.DB.Where("group_id = ? AND tenant_id = ?", groupId, tenantId).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func GetResourceGroupRole(tenantId uint, groupId string, roleId string) (*model.ResourceGroupRole, error) {
	var role model.ResourceGroupRole
	if err := global.DB.Where("id = ? AND group_id = ? AND tenant_id = ?", roleId, groupId, tenantId).
		Find(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func CreateResourceGroupRole(tenantId uint, groupId string, name string, des string, uid string) (*model.ResourceGroupRole, error) {
	if uid == "" {
		uid = utils.GetUuid()
	}
	role := model.ResourceGroupRole{
		Id:          uid,
		TenantId:    tenantId,
		GroupId:     groupId,
		Name:        name,
		Description: des,
	}
	if err := global.DB.Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func UpdateResourceGroupRole(tenantId uint, groupId string, roleId string, name, des string) error {
	if err := global.DB.Model(&model.ResourceGroupRole{}).
		Where("id = ? AND group_id = ? AND tenant_id = ?", roleId, groupId, tenantId).
		Updates(map[string]any{"name": name, "description": des}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroupRole(tenantId uint, groupId string, roleId string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ? AND tenant_id = ?", roleId, tenantId).
			Delete(&model.ResourceGroupRoleAction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND group_id = ? AND tenant_id = ?", roleId, groupId, tenantId).
			Delete(&model.ResourceGroupRole{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetResourceGroupUserRole(tenantId uint, groupId string, userId uint) (*model.ResourceGroupUser, error) {
	var role model.ResourceGroupUser
	if err := global.DB.Where("group_id = ? AND user_id = ? AND tenant_id = ?", groupId, userId, tenantId).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func GetResourceGroupUserActionList(tenantId uint, groupId string, userId uint) ([]model.ResourceGroupRoleAction, error) {
	var userActions []model.ResourceGroupRoleAction
	if err := global.DB.Table("resource_group_role_actions AS rra").
		Select("rra.id", "rra.action_id", "rra.role_id", "rra.tenant_id", "ra.name action_name").
		Joins("LEFT JOIN resource_group_actions AS ra ON ra.id = rra.action_id").
		Joins("LEFT JOIN resource_group_users AS ru ON ru.role_id = rra.role_id").
		Where("ra.group_id = ? AND ru.user_id = ? AND rra.tenant_id", groupId, userId, tenantId).
		Find(&userActions).Error; err != nil {
		return nil, err
	}
	return userActions, nil
}

func GetResourceGroupUserAction(tenantId uint, userId uint, actionId string) (*model.ResourceGroupRoleAction, error) {
	var role model.ResourceGroupRoleAction
	if err := global.DB.Table("resource_group_role_actions as rra").
		Select("rra.id", "rra.action_id", "rra.role_id").
		Joins("LEFT JOIN resource_group_users AS ru ON ru.role_id = rra.role_id").
		Where("rra.action_id = ? AND ru.user_id = ? AND rra.tenant_id = ?", actionId, userId, tenantId).
		First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func CreateResourceGroupUserRole(tenantId uint, groupId string, userId uint, roleId string) (*model.ResourceGroupUser, error) {
	userRole := model.ResourceGroupUser{
		TenantId:     tenantId,
		GroupId:      groupId,
		ClientUserId: userId,
		RoleId:       roleId,
	}
	if err := global.DB.Create(&userRole).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}

func UpdateResourceGroupUserRole(tenantId uint, groupId string, userId uint, roleId string) error {
	if err := global.DB.Model(&model.ResourceGroupUser{}).
		Where("group_id = ? AND user_id = ? AND tenant_id = ?", groupId, userId, tenantId).
		Update("role_id", roleId).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroupUserRole(tenantId uint, groupId string, userId uint) error {
	if err := global.DB.Where("user_id = ? AND group_id = ? AND tenant_id = ?", userId, groupId, tenantId).
		Delete(&model.ResourceGroupUser{}).Error; err != nil {
		return err
	}
	return nil
}

package iam

import (
	"accounts/global"
	"accounts/models"
	"errors"
)

func ListResourceTypeRoleActions(tenantId uint, roleId string) ([]models.ResourceTypeRoleAction, error) {
	var resourceTypeRoleActions []models.ResourceTypeRoleAction
	if err := global.DB.Table("resource_type_role_actions ra").
		Select("ra.id", "ra.role_id", "ra.action_id", "ra.tenant_id", "a.name action_name").
		Joins("LEFT JOIN resource_type_actions a ON ra.action_id = a.id AND ra.tenant_id = a.tenant_id").
		Find(&resourceTypeRoleActions, "role_id = ? AND ra.tenant_id = ?", roleId, tenantId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoleActions, nil
}

func CreateResourceTypeRoleAction(tenantId uint, roleId string, roleAction []models.ResourceTypeRoleAction) error {
	for i := 0; i < len(roleAction); i++ {
		if roleAction[i].ActionId == "" {
			return errors.New("actionId should not be empty")
		}
		roleAction[i].TenantId = tenantId
		roleAction[i].RoleId = roleId
	}
	if err := global.WithTenant(tenantId).Create(roleAction).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceTypeRoleAction(tenantId, roleActionId uint) error {
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, roleActionId).
		Delete(&models.ResourceTypeRoleAction{}).Error; err != nil {
		return err
	}
	return nil
}

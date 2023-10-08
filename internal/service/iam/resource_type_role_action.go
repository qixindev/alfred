package iam

import (
	"alfred/internal/model"
	"alfred/pkg/global"
	"errors"
)

func ListResourceTypeRoleActions(tenantId uint, roleId string) ([]model.ResourceTypeRoleAction, error) {
	var resourceTypeRoleActions []model.ResourceTypeRoleAction
	if err := global.DB.Table("resource_type_role_actions ra").
		Select("ra.id", "ra.role_id", "ra.action_id", "ra.tenant_id", "a.name action_name", "r.name role_name").
		Joins("LEFT JOIN resource_type_actions a ON ra.action_id = a.id").
		Joins("LEFT JOIN resource_type_roles r ON ra.role_id = r.id").
		Find(&resourceTypeRoleActions, "role_id = ? AND ra.tenant_id = ?", roleId, tenantId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoleActions, nil
}

func CreateResourceTypeRoleAction(tenantId uint, roleId string, roleAction []model.ResourceTypeRoleAction) error {
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
		Delete(&model.ResourceTypeRoleAction{}).Error; err != nil {
		return err
	}
	return nil
}

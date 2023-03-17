package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResourceTypeRoleActions(tenantId, roleId uint) ([]models.ResourceTypeRoleAction, error) {
	var resourceTypeRoleActions []models.ResourceTypeRoleAction
	if err := data.DB.Table("resource_type_role_actions ra").
		Select("ra.id", "ra.role_id", "ra.action_id", "ra.tenant_id", "a.name action_name").
		Joins("LEFT JOIN resource_type_actions a ON ra.action_id = a.id AND ra.tenant_id = a.tenant_id").
		Find(&resourceTypeRoleActions, "role_id = ? AND ra.tenant_id = ?", roleId, tenantId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoleActions, nil
}

func GetResourceTypeRoleAction(tenantId, roleActionId uint) (*models.ResourceTypeRoleAction, error) {
	var resourceTypeRoleAction models.ResourceTypeRoleAction
	if err := data.WithTenant(tenantId).Take(&resourceTypeRoleAction, "type_id = ? AND id = ?", roleActionId).Error; err != nil {
		return nil, err
	}
	return &resourceTypeRoleAction, nil
}

func CreateResourceTypeRoleAction(tenantId, roleId uint, roleAction *models.ResourceTypeRoleAction) (*models.ResourceTypeRoleAction, error) {
	roleAction.TenantId = tenantId
	roleAction.RoleId = roleId
	if err := data.WithTenant(tenantId).Create(roleAction).Error; err != nil {
		return nil, err
	}
	return roleAction, nil
}

func UpdateResourceTypeRoleAction(tenantId, roleActionId uint, roleAction *models.ResourceTypeRoleAction) (*models.ResourceTypeRoleAction, error) {
	roleAction.TenantId = tenantId
	roleAction.Id = roleActionId
	if err := data.WithTenant(tenantId).Save(roleAction).Error; err != nil {
		return nil, err
	}
	return roleAction, nil
}

func DeleteResourceTypeRoleAction(tenantId, roleActionId uint) error {
	if err := data.WithTenant(tenantId).Delete(&models.ResourceTypeRoleAction{}, roleActionId).Error; err != nil {
		return err
	}
	return nil
}

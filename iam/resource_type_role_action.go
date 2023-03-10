package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResourceTypeRoleActions(tenantId, typeId, roleId uint) ([]models.ResourceTypeRoleAction, error) {
	var resourceTypeRoleActions []models.ResourceTypeRoleAction
	if err := data.WithTenant(tenantId).Find(&resourceTypeRoleActions, "type_id = ? AND role_id = ?", typeId, roleId).Error; err != nil {
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

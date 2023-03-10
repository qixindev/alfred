package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResourceTypeRoles(tenantId, typeId uint) ([]models.ResourceTypeRole, error) {
	var resourceTypeRoles []models.ResourceTypeRole
	if err := data.WithTenant(tenantId).Find(&resourceTypeRoles, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoles, nil
}

func GetResourceTypeRole(tenantId uint, roleId uint) (*models.ResourceTypeRole, error) {
	var resourceTypeRole models.ResourceTypeRole
	if err := data.WithTenant(tenantId).Take(&resourceTypeRole, "role_id = ?", roleId).Error; err != nil {
		return nil, err
	}
	return &resourceTypeRole, nil
}

func CreateResourceTypeRole(tenantId, typeId uint, role *models.ResourceTypeRole) (*models.ResourceTypeRole, error) {
	role.TenantId = tenantId
	role.TypeId = typeId
	if err := data.WithTenant(tenantId).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func UpdateResourceTypeRole(tenantId, roleId uint, role *models.ResourceTypeRole) (*models.ResourceTypeRole, error) {
	role.TenantId = tenantId
	role.Id = roleId
	if err := data.WithTenant(tenantId).Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func DeleteResourceTypeRole(tenantId, roleId uint) error {
	if err := data.WithTenant(tenantId).Delete(&models.ResourceTypeRole{}, roleId).Error; err != nil {
		return err
	}
	return nil
}

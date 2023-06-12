package iam

import (
	"accounts/global"
	"accounts/models"
)

func ListResourceTypeRoles(tenantId uint, typeId string) ([]models.ResourceTypeRole, error) {
	var resourceTypeRoles []models.ResourceTypeRole
	if err := global.WithTenant(tenantId).Find(&resourceTypeRoles, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoles, nil
}

func CreateResourceTypeRole(tenantId uint, typeId string, role *models.ResourceTypeRole) (*models.ResourceTypeRole, error) {
	role.TenantId = tenantId
	role.TypeId = typeId
	if err := global.WithTenant(tenantId).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func DeleteResourceTypeRole(tenantId uint, roleId string) error {
	if err := global.WithTenant(tenantId).Delete(&models.ResourceTypeRole{}, roleId).Error; err != nil {
		return err
	}
	return nil
}

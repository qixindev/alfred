package iam

import (
	"accounts/global"
	"accounts/models"
)

func ListResourcesRoleUsers(tenantId, resourceId, roleId uint) ([]models.ResourceRoleUser, error) {
	var resourceRoleUsers []models.ResourceRoleUser
	if err := global.WithTenant(tenantId).Find(&resourceRoleUsers, "resource_id = ? AND role_id = ?", resourceId, roleId).Error; err != nil {
		return nil, err
	}
	return resourceRoleUsers, nil
}

func GetResourceRoleUser(tenantId, roleUserId uint) (*models.ResourceRoleUser, error) {
	var resourceRoleUser models.ResourceRoleUser
	if err := global.WithTenant(tenantId).Take(&resourceRoleUser, "id = ?", roleUserId).Error; err != nil {
		return nil, err
	}
	return &resourceRoleUser, nil
}

func CreateResourceRoleUser(tenantId uint, roleUser *models.ResourceRoleUser) (*models.ResourceRoleUser, error) {
	if err := global.WithTenant(tenantId).Create(roleUser).Error; err != nil {
		return nil, err
	}
	return roleUser, nil
}

func UpdateResourceRoleUser(tenantId, roleUserId uint, roleUser *models.ResourceRoleUser) (*models.ResourceRoleUser, error) {
	roleUser.TenantId = tenantId
	roleUser.Id = roleUserId
	if err := global.WithTenant(tenantId).Save(roleUser).Error; err != nil {
		return nil, err
	}
	return roleUser, nil
}

func DeleteResourceRoleUser(tenantId, roleUserId uint) error {
	if err := global.WithTenant(tenantId).Delete(&models.ResourceRoleUser{}, roleUserId).Error; err != nil {
		return err
	}
	return nil
}

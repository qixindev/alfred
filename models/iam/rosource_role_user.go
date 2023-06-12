package iam

import (
	"accounts/global"
	"accounts/models"
)

func ListResourcesRoleUsers(tenantId, resourceId, roleId uint) ([]models.ResourceRoleUser, error) {
	var resourceRoleUsers []models.ResourceRoleUser
	if err := global.DB.Table("resource_role_users as rru").
		Select("rru.id, r.name resource_name, rr.name role_name, rru.client_user_id, cu.sub, u.display_name").
		Joins("LEFT JOIN client_users as cu ON cu.id = rru.client_user_id").
		Joins("LEFT JOIN users as u ON u.id = cu.user_id").
		Joins("LEFT JOIN resources as r ON r.id = rru.resource_id").
		Joins("LEFT JOIN resource_type_roles as rr ON rr.id = rru.role_id").
		Find(&resourceRoleUsers, "rru.tenant_id = ? AND rru.resource_id = ? AND rru.role_id = ?", tenantId, resourceId, roleId).
		Error; err != nil {
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

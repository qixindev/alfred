package iam

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
)

func ListResourcesRoleUsers(tenantId uint, resourceId string, roleId string) ([]model.ResourceRoleUser, error) {
	var resourceRoleUsers []model.ResourceRoleUser
	if err := global.DB.Table("resource_role_users as rru").
		Select("rru.id", "rru.resource_id", "rru.role_id", "r.name resource_name",
			"rr.name role_name", "rru.client_user_id", "cu.sub", "u.display_name").
		Joins("LEFT JOIN client_users as cu ON cu.id = rru.client_user_id").
		Joins("LEFT JOIN users as u ON u.id = cu.user_id").
		Joins("LEFT JOIN resources as r ON r.id = rru.resource_id").
		Joins("LEFT JOIN resource_type_roles as rr ON rr.id = rru.role_id").
		Find(&resourceRoleUsers, "rru.tenant_id = ? AND rru.resource_id = ? AND rru.role_id = ?",
			tenantId, resourceId, roleId).Error; err != nil {
		return nil, err
	}
	return resourceRoleUsers, nil
}

func CreateResourceRoleUser(tenantId uint, roleUser []model.ResourceRoleUser) error {
	if err := global.WithTenant(tenantId).Create(roleUser).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceRoleUser(tenantId, roleUserId uint) error {
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, roleUserId).
		Delete(&model.ResourceRoleUser{}).Error; err != nil {
		return err
	}
	return nil
}

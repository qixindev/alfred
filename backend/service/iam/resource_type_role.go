package iam

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"github.com/google/uuid"
)

func ListResourceTypeRoles(tenantId uint, typeId string) ([]model.ResourceTypeRole, error) {
	var resourceTypeRoles []model.ResourceTypeRole
	if err := global.WithTenant(tenantId).Find(&resourceTypeRoles, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeRoles, nil
}

func CreateResourceTypeRole(tenantId uint, typeId string, role *model.ResourceTypeRole) (*model.ResourceTypeRole, error) {
	role.TenantId = tenantId
	role.TypeId = typeId
	role.Id = uuid.NewString()
	if err := global.WithTenant(tenantId).Create(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func DeleteResourceTypeRole(tenantId uint, roleId string) error {
	if err := global.DB.Where("tenant_id = ? AND role_id = ?", tenantId, roleId).
		Delete(&model.ResourceTypeRoleAction{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND role_id = ?", tenantId, roleId).
		Delete(&model.ResourceRoleUser{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, roleId).
		Delete(&model.ResourceTypeRole{}).Error; err != nil {
		return err
	}
	return nil
}

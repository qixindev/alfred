package iam

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"github.com/google/uuid"
)

func ListResourceTypes(tenantId uint, clientId string) ([]model.ResourceType, error) {
	var resourceTypes []model.ResourceType
	if err := global.WithTenant(tenantId).Find(&resourceTypes, "client_id = ?", clientId).Error; err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

func CreateResourceType(tenantId uint, clientId string, resourceType model.ResourceType) (*model.ResourceType, error) {
	resourceType.TenantId = tenantId
	resourceType.ClientId = clientId
	resourceType.Id = uuid.NewString()
	if err := global.WithTenant(tenantId).Create(&resourceType).Error; err != nil {
		return nil, err
	}
	return &resourceType, nil
}

func DeleteResourceType(tenantId uint, typeId string) error {
	var roles []string
	if err := global.DB.Model(model.ResourceTypeRole{}).Select("id").
		Where("tenant_id = ? and type_id = ?", tenantId, typeId).
		Find(&roles).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND role_id IN ?", tenantId, roles).
		Delete(&model.ResourceTypeRoleAction{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND role_id IN ?", tenantId, roles).
		Delete(&model.ResourceRoleUser{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND type_id = ?", tenantId, typeId).
		Delete(&model.ResourceTypeRole{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND type_id = ?", tenantId, typeId).
		Delete(&model.ResourceTypeAction{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND type_id = ?", tenantId, typeId).
		Delete(&model.Resource{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, typeId).
		Delete(&model.ResourceType{}).Error; err != nil {
		return err
	}
	return nil
}

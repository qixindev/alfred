package iam

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"github.com/google/uuid"
)

func ListResources(tenantId uint, typeId string) ([]model.Resource, error) {
	var resources []model.Resource
	if err := global.WithTenant(tenantId).Find(&resources, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func CreateResource(tenantId uint, typeId string, resource *model.Resource) (*model.Resource, error) {
	resource.TenantId = tenantId
	resource.TypeId = typeId
	resource.Id = uuid.NewString()
	if err := global.WithTenant(tenantId).Create(&resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func UpdateResource(tenantId uint, resource *model.Resource) (*model.Resource, error) {
	if err := global.DB.Model(model.Resource{}).
		Where("id = ? AND tenant_id = ? AND type_id = ?", resource.Id, tenantId, resource.TypeId).
		Update("name", resource.Name).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func DeleteResource(tenantId uint, resourceId string) error {
	if err := global.DB.Where("tenant_id = ? AND resource_id = ?", tenantId, resourceId).
		Delete(&model.ResourceRoleUser{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, resourceId).
		Delete(&model.Resource{}).Error; err != nil {
		return err
	}
	return nil
}

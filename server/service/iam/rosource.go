package iam

import (
	"accounts/global"
	"accounts/models"
	"github.com/google/uuid"
)

func ListResources(tenantId uint, typeId string) ([]models.Resource, error) {
	var resources []models.Resource
	if err := global.WithTenant(tenantId).Find(&resources, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func CreateResource(tenantId uint, typeId string, resource *models.Resource) (*models.Resource, error) {
	resource.TenantId = tenantId
	resource.TypeId = typeId
	resource.Id = uuid.NewString()
	if err := global.WithTenant(tenantId).Create(&resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func UpdateResource(tenantId uint, resource *models.Resource) (*models.Resource, error) {
	if err := global.DB.Model(models.Resource{}).
		Where("id = ? AND tenant_id = ? AND type_id = ?", resource.Id, tenantId, resource.TypeId).
		Update("name", resource.Name).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func DeleteResource(tenantId uint, resourceId string) error {
	if err := global.DB.Where("tenant_id = ? AND resource_id = ?", tenantId, resourceId).
		Delete(&models.ResourceRoleUser{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, resourceId).
		Delete(&models.Resource{}).Error; err != nil {
		return err
	}
	return nil
}

package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResources(tenantId, typeId uint) ([]models.Resource, error) {
	var resources []models.Resource
	if err := data.WithTenant(tenantId).Find(&resources, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func GetResource(tenantId, resourceId uint) (*models.Resource, error) {
	var resource models.Resource
	if err := data.WithTenant(tenantId).Take(&resource, "action_id = ?", resourceId).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func CreateResource(tenantId, typeId uint, resource *models.Resource) (*models.Resource, error) {
	resource.TenantId = tenantId
	resource.TypeId = typeId
	if err := data.WithTenant(tenantId).Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func UpdateResource(tenantId, resourceId uint, resource *models.Resource) (*models.Resource, error) {
	resource.TenantId = tenantId
	resource.Id = resourceId
	if err := data.WithTenant(tenantId).Save(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func DeleteResource(tenantId, resourceId uint) error {
	if err := data.WithTenant(tenantId).Delete(&models.Resource{}, resourceId).Error; err != nil {
		return err
	}
	return nil
}

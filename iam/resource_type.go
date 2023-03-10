package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResourceTypes(tenantId uint) ([]models.ResourceType, error) {
	var resourceTypes []models.ResourceType
	if err := data.WithTenant(tenantId).Find(&resourceTypes).Error; err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

func GetResourceType(tenantId uint, typeId uint) (*models.ResourceType, error) {
	var resourceType models.ResourceType
	if err := data.WithTenant(tenantId).Take(&resourceType, "id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return &resourceType, nil
}

func CreateResourceType(tenantId uint, resourceType *models.ResourceType) (*models.ResourceType, error) {
	resourceType.TenantId = tenantId
	if err := data.WithTenant(tenantId).Create(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func UpdateResourceType(tenantId, typeId uint, resourceType *models.ResourceType) (*models.ResourceType, error) {
	resourceType.Id = typeId
	if err := data.WithTenant(tenantId).Save(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func DeleteResourceType(tenantId, typeId uint) error {
	if err := data.WithTenant(tenantId).Delete(&models.ResourceType{}, typeId).Error; err != nil {
		return err
	}
	return nil
}

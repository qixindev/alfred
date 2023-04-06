package iam

import (
	"accounts/global"
	"accounts/models"
)

func ListResourceTypes(tenantId uint, clientId string) ([]models.ResourceType, error) {
	var resourceTypes []models.ResourceType
	if err := global.WithTenant(tenantId).Find(&resourceTypes, "client_id = ?", clientId).Error; err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

func GetResourceType(tenantId, typeId uint) (*models.ResourceType, error) {
	var resourceType models.ResourceType
	if err := global.WithTenant(tenantId).Take(&resourceType, "id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return &resourceType, nil
}

func CreateResourceType(tenantId uint, clientId string, resourceType *models.ResourceType) (*models.ResourceType, error) {
	resourceType.TenantId = tenantId
	resourceType.ClientId = clientId
	if err := global.WithTenant(tenantId).Create(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func UpdateResourceType(tenantId, typeId uint, resourceType *models.ResourceType) (*models.ResourceType, error) {
	resourceType.Id = typeId
	if err := global.WithTenant(tenantId).Save(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func DeleteResourceType(tenantId, typeId uint) error {
	if err := global.WithTenant(tenantId).Delete(&models.ResourceType{}, typeId).Error; err != nil {
		return err
	}
	return nil
}

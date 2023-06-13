package iam

import (
	"accounts/global"
	"accounts/models"
	"github.com/google/uuid"
)

func ListResourceTypes(tenantId uint, clientId string) ([]models.ResourceType, error) {
	var resourceTypes []models.ResourceType
	if err := global.WithTenant(tenantId).Find(&resourceTypes, "client_id = ?", clientId).Error; err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

func CreateResourceType(tenantId uint, clientId string, resourceType *models.ResourceType) (*models.ResourceType, error) {
	resourceType.TenantId = tenantId
	resourceType.ClientId = clientId
	resourceType.Id = uuid.NewString()
	if err := global.WithTenant(tenantId).Create(resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func DeleteResourceType(tenantId uint, typeId string) error {
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, typeId).
		Delete(&models.ResourceType{}).Error; err != nil {
		return err
	}
	return nil
}

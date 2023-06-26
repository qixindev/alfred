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
	if err := global.WithTenant(tenantId).Create(&resourceType).Error; err != nil {
		return nil, err
	}
	return resourceType, nil
}

func DeleteResourceType(tenantId uint, typeId string) error {
	var roles []uint
	if err := global.DB.Select("id").
		Where("tenant_id = ? and type_id = ?", tenantId, typeId).
		Find(&roles).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND role_id IN ?", tenantId, roles).
		Delete(&models.ResourceTypeRoleAction{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND role_id IN ?", tenantId, roles).
		Delete(&models.ResourceRoleUser{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id IN ?", tenantId, roles).
		Delete(&models.ResourceTypeRole{}).Error; err != nil {
		return err
	}

	if err := global.DB.Where("tenant_id = ? AND type_id = ?", tenantId, typeId).
		Delete(&models.ResourceTypeAction{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND type_id = ?", tenantId, typeId).
		Delete(&models.Resource{}).Error; err != nil {
		return err
	}
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, typeId).
		Delete(&models.ResourceType{}).Error; err != nil {
		return err
	}
	return nil
}

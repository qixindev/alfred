package iam

import (
	"accounts/global"
	"accounts/models"
	"github.com/google/uuid"
)

func ListResourceTypeActions(tenantId uint, typeId string) ([]models.ResourceTypeAction, error) {
	var resourceTypeActions []models.ResourceTypeAction
	if err := global.WithTenant(tenantId).Find(&resourceTypeActions, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeActions, nil
}

func CreateResourceTypeAction(tenantId uint, typeId string, action []models.ResourceTypeAction) error {
	for i := 0; i < len(action); i++ {
		action[i].TenantId = tenantId
		action[i].TypeId = typeId
		action[i].Id = uuid.NewString()
	}
	if err := global.WithTenant(tenantId).Create(action).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceTypeAction(tenantId uint, actionId string) error {
	if err := global.DB.Where("tenant_id = ? AND id = ?", tenantId, actionId).
		Delete(&models.ResourceTypeAction{}).Error; err != nil {
		return err
	}
	return nil
}

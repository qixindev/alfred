package iam

import (
	"accounts/data"
	"accounts/models"
)

func ListResourceTypeActions(tenantId, typeId uint) ([]models.ResourceTypeAction, error) {
	var resourceTypeActions []models.ResourceTypeAction
	if err := data.WithTenant(tenantId).Find(&resourceTypeActions, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeActions, nil
}

func GetResourceTypeAction(tenantId, actionId uint) (*models.ResourceTypeAction, error) {
	var resourceTypeAction models.ResourceTypeAction
	if err := data.WithTenant(tenantId).Take(&resourceTypeAction, "action_id = ?", actionId).Error; err != nil {
		return nil, err
	}
	return &resourceTypeAction, nil
}

func CreateResourceTypeAction(tenantId, typeId uint, action *models.ResourceTypeAction) (*models.ResourceTypeAction, error) {
	action.TenantId = tenantId
	action.TypeId = typeId
	if err := data.WithTenant(tenantId).Create(action).Error; err != nil {
		return nil, err
	}
	return action, nil
}

func UpdateResourceTypeAction(tenantId, actionId uint, action *models.ResourceTypeAction) (*models.ResourceTypeAction, error) {
	action.TenantId = tenantId
	action.Id = actionId
	if err := data.WithTenant(tenantId).Save(action).Error; err != nil {
		return nil, err
	}
	return action, nil
}

func DeleteResourceTypeAction(tenantId, actionId uint) error {
	if err := data.WithTenant(tenantId).Delete(&models.ResourceTypeAction{}, actionId).Error; err != nil {
		return err
	}
	return nil
}

package iam

import (
	"accounts/global"
	"accounts/models"
)

func ListResourceTypeActions(tenantId, typeId uint) ([]models.ResourceTypeAction, error) {
	var resourceTypeActions []models.ResourceTypeAction
	if err := global.WithTenant(tenantId).Find(&resourceTypeActions, "type_id = ?", typeId).Error; err != nil {
		return nil, err
	}
	return resourceTypeActions, nil
}

func GetResourceTypeAction(tenantId, actionId uint) (*models.ResourceTypeAction, error) {
	var resourceTypeAction models.ResourceTypeAction
	if err := global.WithTenant(tenantId).Take(&resourceTypeAction, "action_id = ?", actionId).Error; err != nil {
		return nil, err
	}
	return &resourceTypeAction, nil
}

func CreateResourceTypeAction(tenantId, typeId uint, action []models.ResourceTypeAction) error {
	for i := 0; i < len(action); i++ {
		action[i].TenantId = tenantId
		action[i].TypeId = typeId
	}
	if err := global.WithTenant(tenantId).Create(action).Error; err != nil {
		return err
	}
	return nil
}

func UpdateResourceTypeAction(tenantId, actionId uint, action *models.ResourceTypeAction) (*models.ResourceTypeAction, error) {
	action.TenantId = tenantId
	action.Id = actionId
	if err := global.WithTenant(tenantId).Save(action).Error; err != nil {
		return nil, err
	}
	return action, nil
}

func DeleteResourceTypeAction(tenantId, actionId uint) error {
	if err := global.WithTenant(tenantId).Delete(&models.ResourceTypeAction{}, actionId).Error; err != nil {
		return err
	}
	return nil
}

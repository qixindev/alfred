package rg

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"gorm.io/gorm"
)

func GetResourceGroupActionList(tenantId uint, groupId string) ([]model.ResourceGroupAction, error) {
	var actionsList []model.ResourceGroupAction
	if err := global.DB.Where("group_id = ? AND tenant_id = ?", groupId, tenantId).Find(&actionsList).Error; err != nil {
		return nil, err
	}
	return actionsList, nil
}

func GetResourceGroupAction(tenantId uint, groupId string, actionId string) (*model.ResourceGroupAction, error) {
	var action model.ResourceGroupAction
	if err := global.DB.Where("id = ? AND group_id = ? AND tenant_id = ?", actionId, groupId, tenantId).First(&action).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func CreateResourceGroupAction(tenantId uint, groupId string, name string, des string, uid string) (*model.ResourceGroupAction, error) {
	if uid == "" {
		uid = utils.GetUuid()
	}
	action := model.ResourceGroupAction{
		Id:          uid,
		TenantId:    tenantId,
		GroupId:     groupId,
		Name:        name,
		Description: des,
	}
	if err := global.DB.Create(&action).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func UpdateResourceGroupAction(tenantId uint, groupId string, actionId string, name string) error {
	if err := global.DB.Model(&model.ResourceGroupAction{}).
		Where("id = ? AND group_id = ? AND tenant_id = ?", actionId, groupId, tenantId).
		Update("name", name).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroupAction(tenantId uint, groupId string, actionId string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("action_id = ? AND tenant_id = ?", actionId, tenantId).
			Delete(&model.ResourceGroupRoleAction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? AND group_id = ? AND tenant_id = ?", actionId, groupId, tenantId).
			Delete(&model.ResourceGroupAction{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetResourceGroupRoleActionList(tenantId uint, roleId string) ([]model.ResourceGroupRoleAction, error) {
	var actionList []model.ResourceGroupRoleAction
	if err := global.DB.Where("role_id = ? AND tenant_id = ?", roleId, tenantId).Find(&actionList).Error; err != nil {
		return nil, err
	}
	return actionList, nil
}

func GetResourceGroupRoleAction(tenantId uint, roleId string, actionId string) (*model.ResourceGroupRoleAction, error) {
	var roleAction model.ResourceGroupRoleAction
	if err := global.DB.Where("action_id = ? AND role_id = ? AND tenant_id = ?", actionId, roleId, tenantId).First(&roleAction).Error; err != nil {
		return nil, err
	}
	return &roleAction, nil
}

func CreateResourceGroupRoleAction(tenantId uint, roleId string, actionIds []string) error {
	roleActions := make([]model.ResourceGroupRoleAction, 0)
	for _, actionId := range actionIds {
		roleActions = append(roleActions, model.ResourceGroupRoleAction{
			TenantId: tenantId,
			RoleId:   roleId,
			ActionId: actionId,
		})
	}
	if err := global.DB.Create(&roleActions).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroupRoleAction(tenantId uint, roleId string, actionIds []string) error {
	if err := global.DB.Where("role_id = ? AND tenant_id = ? AND action_id in ?", roleId, tenantId, actionIds).
		Delete(&model.ResourceGroupRoleAction{}).Error; err != nil {
		return err
	}
	return nil
}

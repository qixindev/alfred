package rg

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"gorm.io/gorm"
)

func GetResourceGroupList(tenantId uint, clientId string) ([]model.ResourceGroup, error) {
	var groupList []model.ResourceGroup
	if err := global.DB.Where("tenant_id = ? AND client_id = ?", tenantId, clientId).Find(&groupList).Error; err != nil {
		return nil, err
	}
	return groupList, nil
}

func GetResourceGroup(tenantId uint, clientId string, groupId string) (*model.ResourceGroup, error) {
	var group model.ResourceGroup
	if err := global.DB.Where("id = ? AND client_id = ? AND tenant_id = ?", groupId, clientId, tenantId).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func CreateResourceGroup(tenantId uint, clientId string, name string, uid string) (*model.ResourceGroup, error) {
	if uid == "" {
		uid = utils.GetUuid()
	}
	group := model.ResourceGroup{
		Id:       uid,
		TenantId: tenantId,
		ClientId: clientId,
		Name:     name,
	}
	if err := global.DB.Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func UpdateResourceGroup(tenantId uint, clientId string, groupId string, name string) error {
	if err := global.DB.Where("id = ? AND client_id = ? AND tenant_id = ?", groupId, clientId, tenantId).
		Model(&model.ResourceGroup{}).Update("name", name).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroup(tenantId uint, clientId string, groupId string) error {
	if err := global.DB.Where("id = ? AND client_id = ? AND tenant_id = ?", groupId, clientId, tenantId).
		First(&model.ResourceGroup{}).Error; err != nil {
		return err
	}
	var roles []string
	if err := global.DB.Select("id").Where("group_id = ? AND tenant_id = ?", groupId, clientId, tenantId).
		Model(&model.ResourceGroupRole{}).Find(&roles).Error; err != nil {
		return err
	}
	deleteList := []any{
		model.ResourceGroupUser{},
		model.ResourceGroupAction{},
		model.ResourceGroupRole{},
		model.ResourceGroupResource{},
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id in ? AND tenant_id = ?", roles, tenantId).Delete(&model.ResourceGroupRoleAction{}).Error; err != nil {
			return err
		}
		for _, m := range deleteList {
			if err := tx.Where("group_id = ? AND tenant_id = ?", groupId, tenantId).
				Delete(&m).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("id = ? AND client_id = ? AND tenant_id = ?", groupId, clientId, tenantId).
			Delete(&model.ResourceGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetResourceGroupResourceList(tenantId uint, groupId string) ([]model.ResourceGroupResource, error) {
	var res []model.ResourceGroupResource
	if err := global.DB.Where("group_id = ? AND tenant_id = ?", groupId, tenantId).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func GetResourceGroupResource(tenantId uint, groupId string, resourceId string) (*model.ResourceGroupResource, error) {
	var res model.ResourceGroupResource
	if err := global.DB.Where("id = ? AND group_id = ? AND tenant_id = ?", resourceId, groupId, tenantId).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func CreateResourceGroupResource(tenantId uint, groupId string, name string, uid string) (*model.ResourceGroupResource, error) {
	if uid == "" {
		uid = utils.GetUuid()
	}
	resource := &model.ResourceGroupResource{
		Id:       uid,
		TenantId: tenantId,
		GroupId:  groupId,
		Name:     name,
	}
	if err := global.DB.Create(&resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func UpdateResourceGroupResource(tenantId uint, groupId string, resourceId string, name string) error {
	if err := global.DB.Where("id = ? AND group_id = ? AND tenant_id = ?", resourceId, groupId, tenantId).
		Model(&model.ResourceGroupResource{}).Update("name", name).Error; err != nil {
		return err
	}
	return nil
}

func DeleteResourceGroupResource(tenantId uint, groupId string, resourceId string) error {
	if err := global.DB.Where("id = ? AND group_id = ? AND tenant_id = ?", resourceId, groupId, tenantId).
		Delete(&model.ResourceGroupResource{}).Error; err != nil {
		return err
	}
	return nil
}

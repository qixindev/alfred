package iam

import (
	"accounts/global"
	"accounts/models"
	"errors"
	"fmt"
)

func CheckSinglePermission(tenantId, clientUserId uint, resourceId string, actionId string) (bool, error) {
	var roleActions []models.ResourceTypeRoleAction
	// get all the roles this action supports
	if err := global.DB.Distinct("role_id").Find(&roleActions, "tenant_id = ? AND action_id = ?", tenantId, actionId).Error; err != nil {
		return false, err
	}
	roleIds := make([]string, len(roleActions))
	for i, r := range roleActions {
		roleIds[i] = r.RoleId
	}

	var user models.ResourceRoleUser
	if err := global.DB.First(&user, "tenant_id = ? AND resource_id = ? AND client_user_id = ? AND role_id IN ?", tenantId, resourceId, clientUserId, roleIds); err != nil {
		return false, nil
	}
	return true, nil
}

func CheckPermission(tenantId, clientUserId uint, resourceId string, actionId string) (bool, error) {
	maxDepth := 10
	currentId := resourceId
	for depth := 0; depth < maxDepth; depth++ {
		found, err := CheckSinglePermission(tenantId, clientUserId, currentId, actionId)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}

		var resource models.Resource
		if err = global.DB.First(&resource, "tenant_id = ? AND id = ?", tenantId, currentId).Error; err != nil {
			return false, err
		}

		if resource.ParentId == resource.Id {
			return false, errors.New(fmt.Sprintf("max depth (%d) reached", maxDepth))
		}
		currentId = resource.ParentId
		if currentId == "" {
			break
		}
	}
	return false, nil
}

func GetIamType(tenantId uint, typeId string) (*models.ResourceType, error) {
	var typ models.ResourceType
	if err := global.DB.First(&typ, "tenant_id = ? AND id = ?", tenantId, typeId).Error; err != nil {
		return nil, err
	}
	return &typ, nil
}

func GetIamAction(tenantId uint, typeId, actionId string) (*models.ResourceTypeAction, error) {
	var action models.ResourceTypeAction
	if err := global.DB.First(&action, "tenant_id = ? AND type_id = ? AND id = ?", tenantId, typeId, actionId).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func GetIamRole(tenantId uint, typeId, roleId string) (*models.ResourceTypeRole, error) {
	var role models.ResourceTypeRole
	if err := global.DB.First(&role, "tenant_id = ? AND type_id = ? AND id = ?", tenantId, typeId, roleId).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func GetIamResource(tenantId uint, typeId, resourceId string) (*models.Resource, error) {
	var resource models.Resource
	if err := global.DB.First(&resource, "tenant_id = ? AND type_id = ? AND id = ?", tenantId, typeId, resourceId).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

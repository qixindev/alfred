package iam

import (
	"accounts/data"
	"accounts/models"
	"errors"
	"fmt"
)

func CheckSinglePermission(tenantId, userId, resourceId, actionId uint) (bool, error) {
	var roleActions []models.ResourceTypeRoleAction
	// get all the roles this action supports
	if err := data.DB.Distinct("role_id").Find(&roleActions, "tenant_id = ? AND action_id = ?", tenantId, actionId).Error; err != nil {
		return false, err
	}
	roleIds := make([]uint, len(roleActions))
	for i, r := range roleActions {
		roleIds[i] = r.RoleId
	}

	var user models.ResourceRoleUser
	if err := data.DB.First(&user, "tenant_id = ? AND resource_id = ? AND user_id = ? AND role_id IN ?", tenantId, resourceId, userId, roleIds); err != nil {
		return false, nil
	}
	return true, nil
}

func CheckPermission(tenantId, userId uint, resourceName, actionName string) (bool, error) {
	maxDepth := 10
	var resource models.Resource
	if err := data.DB.First(&resource, "tenant_id = ? AND name = ?", tenantId, resourceName).Error; err != nil {
		return false, err
	}
	var action models.ResourceTypeAction
	if err := data.DB.First(&action, "tenant_id = ? AND name = ?", tenantId, actionName).Error; err != nil {
		return false, err
	}

	resourceId := resource.Id
	for depth := 0; depth < maxDepth; depth++ {
		found, err := CheckSinglePermission(tenantId, userId, resourceId, action.Id)
		if err != nil {
			return false, err
		}
		if resource.ParentId == resource.Id {
			return false, errors.New(fmt.Sprintf("max depth (%d) reached", maxDepth))
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

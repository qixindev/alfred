package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/service/rg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetResourceGroupRoleActionList
// @Summary	获取资源组角色的动作列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [get]
func GetResourceGroupRoleActionList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupRoleActionList(in.Tenant.Id, in.RoleId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupRoleActionList err")
		return
	}
	resp.SuccessWithPaging(c, res, 0)
}

// GetResourceGroupRoleAction
// @Summary	获取资源组角色的动作
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions/{actionsId} [get]
func GetResourceGroupRoleAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupRoleAction(in.Tenant.Id, in.RoleId, in.ActionId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "GetResourceGroupRoleAction err")
		return
	}
	resp.SuccessWithData2(c, res)
}

// CreateResourceGroupRoleAction
// @Summary	创建资源角色的动作
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [post]
func CreateResourceGroupRoleAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.CreateResourceGroupRoleAction(in.Tenant.Id, in.RoleId, in.ActionIds); err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroupRoleAction err")
		return
	}
	resp.Success(c)
}

// UpdateResourceGroupRoleAction
// @Summary	更新资源组角色的动作
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [put]
func UpdateResourceGroupRoleAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	roleActions := make([]model.ResourceGroupRoleAction, 0)
	for _, actionId := range in.ActionIds {
		roleActions = append(roleActions, model.ResourceGroupRoleAction{
			TenantId: in.Tenant.Id,
			RoleId:   in.RoleId,
			ActionId: actionId,
		})
	}
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := global.DB.Where("role_id = ? AND tenant_id = ?", in.RoleId, in.Tenant.Id).
			Delete(&model.ResourceGroupRoleAction{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&roleActions).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroupRoleAction err")
		return
	}

	resp.Success(c)
}

// DeleteResourceGroupRoleAction
// @Summary	删除资源组角色的动作
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [delete]
func DeleteResourceGroupRoleAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroupRoleAction(in.Tenant.Id, in.RoleId, in.ActionIds); err != nil {
		resp.ErrorSqlCreate(c, err, "DeleteResourceGroupRoleAction err")
		return
	}
	resp.Success(c)
}

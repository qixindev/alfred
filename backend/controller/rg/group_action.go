package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/rg"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupActionList
// @Summary	获取资源组角色列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions [get]
func GetResourceGroupActionList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	actions, err := rg.GetResourceGroupActionList(in.Tenant.Id, in.GroupId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupActionList err")
		return
	}
	resp.SuccessWithData(c, actions)
}

// GetResourceGroupAction
// @Summary	获取资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions/{actionId} [get]
func GetResourceGroupAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	action, err := rg.GetResourceGroupAction(in.Tenant.Id, in.GroupId, in.ActionId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupAction err")
		return
	}
	resp.SuccessWithData(c, action)
}

// CreateResourceGroupAction
// @Summary	创建资源角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	role		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions [post]
func CreateResourceGroupAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	action, err := rg.CreateResourceGroupAction(in.Tenant.Id, in.GroupId, in.Name, in.Description, in.Uid)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroupAction err")
		return
	}
	resp.SuccessWithData(c, action)
}

// UpdateResourceGroupAction
// @Summary	更新资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	actionId	path	string		true	"action id"
// @Param	role		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions/{actionId} [put]
func UpdateResourceGroupAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.UpdateResourceGroupAction(in.Tenant.Id, in.GroupId, in.ActionId, in.Name); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroupAction err")
		return
	}
	resp.Success(c)
}

// DeleteResourceGroupAction
// @Summary	删除资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions/{actionId} [delete]
func DeleteResourceGroupAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroupAction(in.Tenant.Id, in.GroupId, in.ActionId); err != nil {
		resp.ErrorSqlDelete(c, err, "DeleteResourceGroupAction err")
		return
	}
	resp.Success(c)
}

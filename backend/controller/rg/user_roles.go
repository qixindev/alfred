package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/rg"
	"errors"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupUserList
// @Summary	组内用户列表
// @Tags	resource-group
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	groupId		path	string	true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users [get]
func GetResourceGroupUserList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupUserList(in.Tenant.Id, in.GroupId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupUserList err")
		return
	}
	resp.SuccessWithPaging(c, res, 0)
}

// GetResourceGroupUserRole
// @Summary	用户在组内的角色
// @Tags	resource-group
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	groupId		path	string	true	"group id"
// @Param	userId		path	integer	true	"client user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/roles [get]
func GetResourceGroupUserRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupUserRole(in.Tenant.Id, in.GroupId, in.UserId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupUserRole err")
		return
	}
	resp.SuccessWithData2(c, res)
}

// GetResourceGroupUserActionList
// @Summary	用户在组内所拥有的权限列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	integer		true	"client user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/actions [get]
func GetResourceGroupUserActionList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupUserActionList(in.Tenant.Id, in.GroupId, in.UserId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupUserActionList err")
		return
	}
	resp.SuccessWithPaging(c, res, 0)
}

// GetResourceGroupUserAction
// @Summary	用户在组内是否拥有某个权限
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	string		true	"client user id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/actions/{actionId} [get]
func GetResourceGroupUserAction(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupUserAction(in.Tenant.Id, in.UserId, in.ActionId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "GetResourceGroupUserAction err")
		return
	}
	resp.SuccessWithData2(c, res)
}

// CreateResourceGroupUserRole
// @Summary	将用户拉入组内
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	integer		true	"client user id"
// @Param	group		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [post]
func CreateResourceGroupUserRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if in.RoleId == "" {
		resp.ErrorRequest(c, errors.New("body roleId should not be empty"))
		return
	}
	res, err := rg.CreateResourceGroupUserRole(in.Tenant.Id, in.GroupId, in.UserId, in.RoleId)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroupUserRole err")
		return
	}
	resp.SuccessWithData2(c, res)
}

// UpdateResourceGroupUserRole
// @Summary	修改用户在组内的角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	integer		true	"client user id"
// @Param	group		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [put]
func UpdateResourceGroupUserRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if in.RoleId == "" {
		resp.ErrorRequest(c, errors.New("body roleId should not be empty"))
		return
	}
	if err := rg.UpdateResourceGroupUserRole(in.Tenant.Id, in.GroupId, in.UserId, in.RoleId); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroupUserRole err")
		return
	}
	resp.Success(c)
}

// DeleteResourceGroupUser
// @Summary	踢出用户
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	integer		true	"client user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [delete]
func DeleteResourceGroupUser(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroupUserRole(in.Tenant.Id, in.GroupId, in.UserId); err != nil {
		resp.ErrorSqlDelete(c, err, "DeleteResourceGroupUserRole err")
		return
	}
	resp.Success(c)
}

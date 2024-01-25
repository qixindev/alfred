package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/rg"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupRoleList
// @Summary	获取资源组角色列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles [get]
func GetResourceGroupRoleList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupRoleList(in.Tenant.Id, in.GroupId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupRoleList err")
		return
	}
	resp.SuccessWithData(c, res)
}

// GetResourceGroupRole
// @Summary	获取资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId} [get]
func GetResourceGroupRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupRole(in.Tenant.Id, in.GroupId, in.RoleId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "GetResourceGroupRole err")
		return
	}
	resp.SuccessWithData(c, res)
}

// CreateResourceGroupRole
// @Summary	创建资源角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles [post]
func CreateResourceGroupRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.CreateResourceGroupRole(in.Tenant.Id, in.GroupId, in.Name, in.Description, in.Uid)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroupRole err")
		return
	}
	resp.SuccessWithData(c, res)
}

// UpdateResourceGroupRole
// @Summary	更新资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId} [put]
func UpdateResourceGroupRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.UpdateResourceGroupRole(in.Tenant.Id, in.GroupId, in.RoleId, in.Name); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroupRole err")
		return
	}
	resp.Success(c)
}

// DeleteResourceGroupRole
// @Summary	删除资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId} [delete]
func DeleteResourceGroupRole(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroupRole(in.Tenant.Id, in.GroupId, in.RoleId); err != nil {
		resp.ErrorSqlUpdate(c, err, "DeleteResourceGroupRole err")
		return
	}
	resp.Success(c)
}

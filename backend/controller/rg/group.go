package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/rg"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupList
// @Summary	获取资源组列表
// @Tags	resource-group
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups [get]
func GetResourceGroupList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	groups, err := rg.GetResourceGroupList(in.Tenant.Id, in.ClientId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupList err")
		return
	}
	resp.SuccessWithData(c, groups)
}

// GetResourceGroup
// @Summary	获取资源组详细信息
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId} [get]
func GetResourceGroup(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	group, err := rg.GetResourceGroup(in.Tenant.Id, in.ClientId, in.GroupId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "GetResourceGroup err")
		return
	}
	resp.SuccessWithData(c, group)
}

// CreateResourceGroup
// @Summary	创建资源组
// @Tags	resource-group
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	client		path	string				true	"client"	default(default)
// @Param	group		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups [post]
func CreateResourceGroup(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	group, err := rg.CreateResourceGroup(in.Tenant.Id, in.ClientId, in.Name, in.Description, in.Uid)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroup err")
		return
	}
	resp.SuccessWithData(c, group)
}

// UpdateResourceGroup
// @Summary	更新资源组
// @Tags	resource-group
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	client		path	string				true	"client"	default(default)
// @Param	groupId		path	string				true	"group id"
// @Param	group		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId} [put]
func UpdateResourceGroup(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.UpdateResourceGroup(in.Tenant.Id, in.ClientId, in.GroupId, in.Name, in.Description); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroup err")
		return
	}
	resp.Success(c)
}

// DeleteResourceGroup
// @Summary	删除资源组
// @Tags	resource-group
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	client		path	string				true	"client"	default(default)
// @Param	groupId		path	string				true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId} [delete]
func DeleteResourceGroup(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroup(in.Tenant.Id, in.ClientId, in.GroupId); err != nil {
		resp.ErrorSqlDelete(c, err, "DeleteResourceGroup err")
		return
	}
	resp.Success(c)
}

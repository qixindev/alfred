package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/rg"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupResourceList
// @Summary	获取资源组的资源列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources [get]
func GetResourceGroupResourceList(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupResourceList(in.Tenant.Id, in.GroupId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "GetResourceGroupResourceList err")
		return
	}
	resp.SuccessWithData(c, res)
}

// GetResourceGroupResource
// @Summary	获取资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [get]
func GetResourceGroupResource(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.GetResourceGroupResource(in.Tenant.Id, in.GroupId, in.ResourceId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "GetResourceGroupResource err")
		return
	}
	resp.SuccessWithData(c, res)
}

// CreateResourceGroupResource
// @Summary	创建资源的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources [post]
func CreateResourceGroupResource(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	res, err := rg.CreateResourceGroupResource(in.Tenant.Id, in.GroupId, in.Name, in.Description, in.Uid)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "CreateResourceGroupResource err")
		return
	}
	resp.SuccessWithData(c, res)
}

// UpdateResourceGroupResource
// @Summary	更新资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Param	data		body	model.RequestResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [put]
func UpdateResourceGroupResource(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUriAndJson(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.UpdateResourceGroupResource(in.Tenant.Id, in.GroupId, in.ResourceId, in.Name, in.Description); err != nil {
		resp.ErrorSqlUpdate(c, err, "UpdateResourceGroupResource err")
		return
	}
	resp.Success(c)
}

// DeleteResourceGroupResource
// @Summary	删除资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [delete]
func DeleteResourceGroupResource(c *gin.Context) {
	var in model.RequestResourceGroup
	if err := internal.BindUri(c, &in).SetTenant(&in.Tenant).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if err := rg.DeleteResourceGroupResource(in.Tenant.Id, in.GroupId, in.ResourceId); err != nil {
		resp.ErrorSqlDelete(c, err, "DeleteResourceGroupResource err")
		return
	}
	resp.Success(c)
}

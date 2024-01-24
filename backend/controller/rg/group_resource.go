package rg

import (
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
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

}

// CreateResourceGroupResource
// @Summary	创建资源的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	role		body	model.ResourceGroupResource	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources [post]
func CreateResourceGroupResource(c *gin.Context) {
	var in model.ResourceGroupResource
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrorRequest(c, err)
	}

}

// UpdateResourceGroupResource
// @Summary	更新资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Param	role		body	model.ResourceGroupResource	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [put]
func UpdateResourceGroupResource(c *gin.Context) {

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

}

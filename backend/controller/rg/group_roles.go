package rg

import (
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
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

}

// CreateResourceGroupRole
// @Summary	创建资源角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	role		body	model.ResourceGroupRole	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles [post]
func CreateResourceGroupRole(c *gin.Context) {
	var in model.ResourceGroupRole
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrorRequest(c, err)
	}

}

// UpdateResourceGroupRole
// @Summary	更新资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	role		body	model.ResourceGroupRole	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId} [put]
func UpdateResourceGroupRole(c *gin.Context) {

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

}

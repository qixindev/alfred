package rg

import "github.com/gin-gonic/gin"

// GetResourceGroupRoleActionList
// @Summary	获取资源组角色列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [get]
func GetResourceGroupRoleActionList(c *gin.Context) {

}

// GetResourceGroupRoleAction
// @Summary	获取资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions/{actionsId} [get]
func GetResourceGroupRoleAction(c *gin.Context) {

}

// CreateResourceGroupRoleAction
// @Summary	创建资源角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	role		body	model.ResourceGroupRoleAction	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [post]
func CreateResourceGroupRoleAction(c *gin.Context) {

}

// UpdateResourceGroupRoleAction
// @Summary	更新资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Param	role		body	model.ResourceGroupRoleAction	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [put]
func UpdateResourceGroupRoleAction(c *gin.Context) {

}

// DeleteResourceGroupRoleAction
// @Summary	删除资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	roleId		path	string		true	"role id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/roles/{roleId}/actions [put]
func DeleteResourceGroupRoleAction(c *gin.Context) {

}

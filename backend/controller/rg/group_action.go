package rg

import "github.com/gin-gonic/gin"

// GetResourceGroupActionList
// @Summary	获取资源组角色列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions [get]
func GetResourceGroupActionList(c *gin.Context) {

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

}

// CreateResourceGroupAction
// @Summary	创建资源角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	role		body	model.ResourceGroupAction	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions [post]
func CreateResourceGroupAction(c *gin.Context) {

}

// UpdateResourceGroupAction
// @Summary	更新资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	actionId	path	string		true	"action id"
// @Param	role		body	model.ResourceGroupAction	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions/{actionId} [put]
func UpdateResourceGroupAction(c *gin.Context) {

}

// DeleteResourceGroupAction
// @Summary	删除资源组角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/actions/{actionId} [put]
func DeleteResourceGroupAction(c *gin.Context) {

}

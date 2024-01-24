package rg

import (
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"github.com/gin-gonic/gin"
)

// GetResourceGroupUserRole
// @Summary	用户在组内的角色
// @Tags	resource-group
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	groupId		path	string	true	"group id"
// @Param	userId		path	string	true	"user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/roles [get]
func GetResourceGroupUserRole(c *gin.Context) {

}

// GetResourceGroupUserActionList
// @Summary	用户在组内所拥有的权限列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	string		true	"user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/actions [get]
func GetResourceGroupUserActionList(c *gin.Context) {

}

// GetResourceGroupUserAction
// @Summary	用户在组内是否拥有某个权限
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	string		true	"user id"
// @Param	actionId	path	string		true	"action id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId}/actions/{actionId} [get]
func GetResourceGroupUserAction(c *gin.Context) {

}

// CreateResourceGroupUserRole
// @Summary	将用户拉入组内
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	userId		path	string		true	"user id"
// @Param	group		body	model.ResourceGroupUser	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [post]
func CreateResourceGroupUserRole(c *gin.Context) {
	var in model.ResourceGroupUser
	if err := c.ShouldBindJSON(&in); err != nil {
		resp.ErrorRequest(c, err)
	}

}

// UpdateResourceGroupUserRole
// @Summary	修改用户在组内的角色
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	string		true	"user id"
// @Param	group		body	model.ResourceGroupUser	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [put]
func UpdateResourceGroupUserRole(c *gin.Context) {

}

// DeleteResourceGroupUser
// @Summary	踢出用户
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	userId		path	string		true	"user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/users/{userId} [delete]
func DeleteResourceGroupUser(c *gin.Context) {

}

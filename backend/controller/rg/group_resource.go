package rg

import "github.com/gin-gonic/gin"

// GetGroupResourceList
// @Summary	获取资源组的资源列表
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources [get]
func GetGroupResourceList(c *gin.Context) {

}

// GetGroupResource
// @Summary	获取资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [get]
func GetGroupResource(c *gin.Context) {

}

// CreateGroupResource
// @Summary	创建资源的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	role		body	model.GroupResource	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources [post]
func CreateGroupResource(c *gin.Context) {

}

// UpdateGroupResource
// @Summary	更新资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Param	role		body	model.GroupResource	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [put]
func UpdateGroupResource(c *gin.Context) {

}

// DeleteGroupResource
// @Summary	删除资源组的资源
// @Tags	resource-group
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	groupId		path	string		true	"group id"
// @Param	resourceId	path	string		true	"resource id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId}/resources/{resourceId} [put]
func DeleteGroupResource(c *gin.Context) {

}

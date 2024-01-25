package rg

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
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
	if err := internal.BindUri(c, &in).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
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

}

// CreateResourceGroup
// @Summary	创建资源组
// @Tags	resource-group
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	client		path	string				true	"client"	default(default)
// @Param	group		body	model.ResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups [post]
func CreateResourceGroup(c *gin.Context) {
	var in model.ResourceGroup
	if err := internal.BindJson(c, &in).Error; err != nil {
		resp.ErrorRequest(c, err)
	}

}

// UpdateResourceGroup
// @Summary	更新资源组
// @Tags	resource-group
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	client		path	string				true	"client"	default(default)
// @Param	groupId		path	string				true	"group id"
// @Param	group		body	model.ResourceGroup	true	"body"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/resourceGroups/{groupId} [put]
func UpdateResourceGroup(c *gin.Context) {

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

}

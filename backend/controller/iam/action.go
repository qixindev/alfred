package iam

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/iam"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamAction
// @Summary	获取操作列表
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions [get]
func ListIamAction(c *gin.Context) {
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	actions, err := iam.ListResourceTypeActions(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "get resource type action err", true)
		return
	}
	resp.SuccessWithArrayData(c, actions, 0)
}

// NewIamAction godoc
//
// @Summary	new iam action
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string		true	"tenant"
// @Param	iamBody		body	[]model.ResourceTypeAction	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions [post]
func NewIamAction(c *gin.Context) {
	var action []model.ResourceTypeAction
	if err := c.BindJSON(&action); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	typ, err := iam.GetIamType(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get iam resource type err")
		return
	}

	if err = iam.CreateResourceTypeAction(tenant.Id, typ.Id, action); err != nil {
		resp.ErrorSqlCreate(c, err, "create resource type action err")
		return
	}
	resp.SuccessWithData(c, action)
}

// DeleteIamAction
// @Summary	delete iam action
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	actionId	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions/{actionId} [delete]
func DeleteIamAction(c *gin.Context) {
	actionId := c.Param("actionId")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	action, err := iam.GetIamAction(tenant.Id, typeId, actionId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get action err")
		return
	}

	if err = iam.DeleteResourceTypeAction(tenant.Id, action.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource type action err")
		return
	}
	c.Status(http.StatusNoContent)
}

// ListIamRoleAction
// @Summary	get iam role action list
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions [get]
func ListIamRoleAction(c *gin.Context) {
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	roleActions, err := iam.ListResourceTypeRoleActions(tenant.Id, roleId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource type role action err", true)
		return
	}
	resp.SuccessWithArrayData(c, roleActions, 0)
}

// NewIamRoleAction
// @Summary	new iam role action
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Param	iamBody		body	[]model.ResourceTypeRoleAction	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions [post]
func NewIamRoleAction(c *gin.Context) {
	var roleAction []model.ResourceTypeRoleAction
	if err := c.BindJSON(&roleAction); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	roleId := c.Param("roleId")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	role, err := iam.GetIamRole(tenant.Id, typeId, roleId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get role err")
		return
	}

	if err = iam.CreateResourceTypeRoleAction(tenant.Id, role.Id, roleAction); err != nil {
		resp.ErrorSqlCreate(c, err, "create role action err")
		return
	}
	resp.SuccessWithData(c, roleAction)
}

// DeleteIamRoleAction
// @Summary	delete iam role action
// @Tags	iam-action
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Param	actionId	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions/{actionId} [delete]
func DeleteIamRoleAction(c *gin.Context) {
	actionId := c.Param("actionId")
	roleId := c.Param("roleId")
	var roleAction model.ResourceTypeRoleAction
	if err := internal.TenantDB(c).First(&roleAction, "role_id = ? AND action_id = ?", roleId, actionId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get role action err")
		return
	}
	if err := iam.DeleteResourceTypeRoleAction(roleAction.TenantId, roleAction.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource role action err")
		return
	}
	c.Status(http.StatusNoContent)
}

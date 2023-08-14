package iam

import (
	"accounts/internal/controller/internal"
	"accounts/internal/model"
	"accounts/internal/service/iam"
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// ListIamAction godoc
//
//	@Summary		iam action
//	@Schemes
//	@Description	get iam action list
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions [get]
func ListIamAction(c *gin.Context) {
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	actions, err := iam.ListResourceTypeActions(tenant.Id, typeId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource type action")
		global.LOG.Error("list resource type action err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, actions)
}

// NewIamAction godoc
//
//	@Summary		iam action
//	@Schemes
//	@Description	new iam action
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string		true	"tenant"
//	@Param			iamBody		body	[]internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions [post]
func NewIamAction(c *gin.Context) {
	var action []model.ResourceTypeAction
	if err := c.BindJSON(&action); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	typ, err := iam.GetIamType(tenant.Id, typeId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such iam resource type")
		global.LOG.Error("get iam type err: " + err.Error())
		return
	}

	if err = iam.CreateResourceTypeAction(tenant.Id, typ.Id, action); err != nil {
		internal.ErrorSqlResponse(c, "failed to create resource type action")
		global.LOG.Error("create resource type action err: " + err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteIamAction godoc
//
//	@Summary		iam action
//	@Schemes
//	@Description	delete iam action
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			actionId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions/{actionId} [delete]
func DeleteIamAction(c *gin.Context) {
	actionId := c.Param("actionId")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	action, err := iam.GetIamAction(tenant.Id, typeId, actionId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such action")
		global.LOG.Error("get iam action err: " + err.Error())
		return
	}

	if err = iam.DeleteResourceTypeAction(tenant.Id, action.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource type action")
		global.LOG.Error("delete resource type action err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// ListIamRoleAction godoc
//
//	@Summary		iam role action
//	@Schemes
//	@Description	get iam role action list
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions [get]
func ListIamRoleAction(c *gin.Context) {
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	roleActions, err := iam.ListResourceTypeRoleActions(tenant.Id, roleId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource type role action list")
		global.LOG.Error("list resource type role action err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, roleActions)
}

// NewIamRoleAction godoc
//
//	@Summary		iam role action
//	@Schemes
//	@Description	new iam role action
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			iamBody		body	[]internal.IamActionRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions [post]
func NewIamRoleAction(c *gin.Context) {
	var roleAction []model.ResourceTypeRoleAction
	if err := c.BindJSON(&roleAction); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	roleId := c.Param("roleId")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	role, err := iam.GetIamRole(tenant.Id, typeId, roleId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such role")
		global.LOG.Error("get iam role err: ", zap.Error(err))
		return
	}

	if err = iam.CreateResourceTypeRoleAction(tenant.Id, role.Id, roleAction); err != nil {
		internal.ErrorSqlResponse(c, "failed to create role action")
		global.LOG.Error("create resource type role action err: " + err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteIamRoleAction godoc
//
//	@Summary		iam role action
//	@Schemes
//	@Description	delete iam role action
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			actionId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/actions/{actionId} [delete]
func DeleteIamRoleAction(c *gin.Context) {
	actionId := c.Param("actionId")
	roleId := c.Param("roleId")
	var roleAction model.ResourceTypeRoleAction
	if err := internal.TenantDB(c).First(&roleAction, "role_id = ? AND action_id = ?", roleId, actionId).Error; err != nil {
		internal.ErrReqParaCustom(c, "no such role action")
		global.LOG.Error("get resource type role action err: " + err.Error())
		return
	}
	if err := iam.DeleteResourceTypeRoleAction(roleAction.TenantId, roleAction.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource role action")
		global.LOG.Error("delete resource type role action err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
	iam2 "accounts/server/service/iam"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamAction godoc
//
//	@Summary		iam action
//	@Schemes
//	@Description	get iam action list
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/actions [get]
func ListIamAction(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get iam type err: " + err.Error())
		return
	}
	actions, err := iam2.ListResourceTypeActions(typ.TenantId, typ.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/actions [post]
func NewIamAction(c *gin.Context) {
	var action []models.ResourceTypeAction
	if err := c.BindJSON(&action); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get iam type err: " + err.Error())
		return
	}
	if err = iam2.CreateResourceTypeAction(typ.TenantId, typ.Id, action); err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("CreateResourceTypeAction err: " + err.Error())
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
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			action		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/actions/{action} [delete]
func DeleteIamAction(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("getType err: " + err.Error())
		return
	}
	actionName := c.Param("action")
	var action models.ResourceTypeAction
	if err = internal.TenantDB(c).First(&action, "type_id = ? AND name = ?", typ.Id, actionName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource type action err: " + err.Error())
		return
	}
	if err = iam2.DeleteResourceTypeAction(typ.TenantId, action.Id); err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions [get]
func ListIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get role err: " + err.Error())
		return
	}
	roleActions, err := iam2.ListResourceTypeRoleActions(role.TenantId, role.Id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions [post]
func NewIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get role")
		global.LOG.Error("get role err: " + err.Error())
		return
	}
	var roleAction []models.ResourceTypeRoleAction
	if err = c.BindJSON(&roleAction); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	if err = iam2.CreateResourceTypeRoleAction(role.TenantId, role.Id, roleAction); err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Param			action		path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions/{action} [delete]
func DeleteIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get role err: " + err.Error())
		return
	}
	actionName := c.Param("action")
	var action models.ResourceTypeAction
	if err = internal.TenantDB(c).First(&action, "type_id = ? AND name = ?", role.TypeId, actionName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource type action err: " + err.Error())
		return
	}
	var roleAction models.ResourceTypeRoleAction
	if err = internal.TenantDB(c).First(&roleAction, "role_id = ? AND action_id = ?", role.Id, action.Id).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource type role action err: " + err.Error())
		return
	}
	if iam2.DeleteResourceTypeRoleAction(role.TenantId, roleAction.Id) != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("delete resource type role action err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

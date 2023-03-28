package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/models/iam"
	"accounts/server/internal"
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
	actions, err := iam.ListResourceTypeActions(typ.TenantId, typ.Id)
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
	var action models.ResourceTypeAction
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
	a, err := iam.CreateResourceTypeAction(typ.TenantId, typ.Id, &action)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("CreateResourceTypeAction err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, a)
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
	if err = iam.DeleteResourceTypeAction(typ.TenantId, action.Id); err != nil {
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
	roleActions, err := iam.ListResourceTypeRoleActions(role.TenantId, role.Id)
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
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get role err: " + err.Error())
		return
	}
	var roleAction models.ResourceTypeRoleAction
	if err = c.BindJSON(&roleAction); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	ra, err := iam.CreateResourceTypeRoleAction(role.TenantId, role.Id, &roleAction)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("create resource type role action err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, ra)
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
	if iam.DeleteResourceTypeRoleAction(role.TenantId, roleAction.Id) != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("delete resource type role action err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetIamActionUser godoc
//
//	@Summary		iam action user
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			action		path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/actions/{action}/users/{user} [get]
func GetIamActionUser(c *gin.Context) {
	resource, err := getResource(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource err: " + err.Error())
		return
	}
	action, err := getAction(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get action err: " + err.Error())
		return
	}

	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}
	userName := c.Param("user")
	var clientUser models.ClientUser
	if err = internal.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", client.Id, userName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client user err: " + err.Error())
		return
	}

	result, err := iam.CheckPermission(resource.TenantId, clientUser.Id, resource.Id, action.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("check permission err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

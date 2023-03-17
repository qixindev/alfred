package iam

import (
	"accounts/iam"
	"accounts/middlewares"
	"accounts/models"
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/actions [get]
func ListIamAction(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	actions, err := iam.ListResourceTypeActions(typ.TenantId, typ.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/actions [post]
func NewIamAction(c *gin.Context) {
	var action models.ResourceTypeAction
	if c.BindJSON(&action) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	a, err := iam.CreateResourceTypeAction(typ.TenantId, typ.Id, &action)
	if err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/actions/{action} [delete]
func DeleteIamAction(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	actionName := c.Param("action")
	var action models.ResourceTypeAction
	if err := middlewares.TenantDB(c).First(&action, "type_id = ? AND name = ?", typ.Id, actionName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := iam.DeleteResourceTypeAction(typ.TenantId, action.Id); err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions [get]
func ListIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	roleActions, err := iam.ListResourceTypeRoleActions(role.TenantId, role.Id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions [post]
func NewIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var roleAction models.ResourceTypeRoleAction
	if c.BindJSON(&roleAction) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	ra, err := iam.CreateResourceTypeRoleAction(role.TenantId, role.Id, &roleAction)
	if err != nil {
		c.Status(http.StatusBadRequest)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles/{role}/actions/{action} [delete]
func DeleteIamRoleAction(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	actionName := c.Param("action")
	var action models.ResourceTypeAction
	if err := middlewares.TenantDB(c).First(&action, "type_id = ? AND name = ?", role.TypeId, actionName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var roleAction models.ResourceTypeRoleAction
	if err := middlewares.TenantDB(c).First(&roleAction, "role_id = ? AND action_id = ?", role.Id, action.Id).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if iam.DeleteResourceTypeRoleAction(role.TenantId, roleAction.Id) != nil {
		c.Status(http.StatusBadRequest)
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
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/actions/{action}/users/{user} [get]
func GetIamActionUser(c *gin.Context) {
	resource, err := getResource(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	action, err := getAction(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userName := c.Param("user")
	var clientUser models.ClientUser
	if err := middlewares.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", client.Id, userName); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := iam.CheckPermission(resource.TenantId, clientUser.Id, resource.Id, action.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, result)
}

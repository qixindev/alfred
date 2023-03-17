package iam

import (
	"accounts/iam"
	"accounts/middlewares"
	"accounts/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamRole godoc
//
//	@Summary	iam role
//	@Schemes
//	@Description	get iam role list
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles [get]
func ListIamRole(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	roles, err := iam.ListResourceTypeRoles(typ.TenantId, typ.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// NewIamRole godoc
//
//	@Summary		iam role
//	@Schemes
//	@Description	new iam role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles [post]
func NewIamRole(c *gin.Context) {
	var role models.ResourceTypeRole
	if c.BindJSON(&role) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	r, err := iam.CreateResourceTypeRole(typ.TenantId, typ.Id, &role)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, r)
}

// DeleteIamRole godoc
//
//	@Summary		iam role
//	@Schemes
//	@Description	delete iam role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/roles/{role} [delete]
func DeleteIamRole(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	roleName := c.Param("role")
	var role models.ResourceTypeRole
	if err := middlewares.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := iam.DeleteResourceTypeRole(typ.TenantId, role.Id); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListIamResourceRole godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	get iam resource role list
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users [get]
func ListIamResourceRole(c *gin.Context) {
	resource, role, err := getResourceAndRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	roleUsers, err := iam.ListResourcesRoleUsers(resource.TenantId, resource.Id, role.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, roleUsers)
}

// NewIamResourceRole godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	new iam resource role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users [post]
func NewIamResourceRole(c *gin.Context) {
	var roleUser models.ResourceRoleUser
	if c.BindJSON(&roleUser) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	resource, role, err := getResourceAndRole(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	ru, err := iam.CreateResourceRoleUser(resource.TenantId, resource.Id, role.Id, &roleUser)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, ru)
}

// DeleteIamResourceRole godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	delete iam resource role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users/{user} [delete]
func DeleteIamResourceRole(c *gin.Context) {
	resource, role, err := getResourceAndRole(c)
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
	var roleUser models.ResourceRoleUser
	if err := middlewares.TenantDB(c).First(&roleUser, "resource_id = ? AND role_id = ? and client_user_id", resource.Id, role.Id, clientUser.Id).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := iam.DeleteResourceRoleUser(resource.TenantId, roleUser.Id); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusNoContent)
}

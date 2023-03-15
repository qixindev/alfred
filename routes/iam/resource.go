package iam

import (
	"accounts/iam"
	"accounts/middlewares"
	"accounts/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamResourceType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	get iam resource type list
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types [get]
func ListIamResourceType(c *gin.Context) {
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	types, err := iam.ListResourceTypes(client.TenantId, client.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, types)
}

// NewIamResourceType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	new iam resource type
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types [post]
func NewIamResourceType(c *gin.Context) {
	var typ models.ResourceType
	if c.BindJSON(&typ) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	t, err := iam.CreateResourceType(client.TenantId, client.Id, &typ)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, t)
}

// DeleteIamResourceType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	delete iam resource type
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type} [delete]
func DeleteIamResourceType(c *gin.Context) {
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	typeName := c.Param("type")
	var typ models.ResourceType
	if err := middlewares.TenantDB(c).First(&typ, "client_id = ? AND name = ?", client.Id, typeName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := iam.DeleteResourceType(client.TenantId, typ.Id); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListIamRole godoc
//
//	@Summary	iam role
//	@Schemes
//	@Description	get iam role list
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles [get]
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles [post]
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			role		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/roles/{role} [delete]
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

// ListIamAction godoc
//
//	@Summary		iam action
//	@Schemes
//	@Description	get iam action list
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/actions [post]
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			action		path	integer	true	"tenant"
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
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
	roleActions, err := iam.ListResourceTypeRoleActions(role.TenantId, role.TypeId, role.Id)
	c.JSON(http.StatusOK, roleActions)
}

// NewIamRoleAction godoc
//
//	@Summary		iam role action
//	@Schemes
//	@Description	new iam role action
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
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

// ListIamResource godoc
//
//	@Summary		iam resource
//	@Schemes
//	@Description	get iam resource list
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources [get]
func ListIamResource(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	resources, err := iam.ListResources(typ.TenantId, typ.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, resources)
}

// NewIamResource godoc
//
//	@Summary		iam resource
//	@Schemes
//	@Description	new iam resource
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources [post]
func NewIamResource(c *gin.Context) {
	var resource models.Resource
	if c.BindJSON(&resource) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	r, err := iam.CreateResource(typ.TenantId, typ.Id, &resource)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, r)
}

// DeleteIamResource godoc
//
//	@Summary		iam resource
//	@Schemes
//	@Description	delete iam resource
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources/{resource} [delete]
func DeleteIamResource(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	resourceName := c.Param("resource")
	var resource models.Resource
	if err := middlewares.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := iam.DeleteResource(typ.TenantId, resource.Id); err != nil {
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users [get]
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users [post]
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
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			resource	path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/roles/{role}/users/{user} [delete]
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

// GetIamActionUser godoc
//
//	@Summary		iam action user
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
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

func AddIamRoutes(rg *gin.RouterGroup) {
	rg.GET("/types", ListIamResourceType)
	rg.POST("/types", NewIamResourceType)
	rg.DELETE("/types/:type", DeleteIamResourceType)

	rg.GET("/types/:type/roles", ListIamRole)
	rg.POST("/types/:type/roles", NewIamRole)
	rg.DELETE("/types/:type/roles/:role", DeleteIamRole)

	rg.GET("/types/:type/actions", ListIamAction)
	rg.POST("/types/:type/actions", NewIamAction)
	rg.DELETE("/types/:type/actions/:action", DeleteIamAction)

	rg.GET("/types/:type/roles/:role/actions", ListIamRoleAction)
	rg.POST("/types/:type/roles/:role/actions", NewIamRoleAction)
	rg.DELETE("/types/:type/roles/:role/actions/:action", DeleteIamRoleAction)

	rg.GET("/types/:type/resources", ListIamResource)
	rg.POST("/types/:type/resources", NewIamResource)
	rg.DELETE("/types/:type/resources/:resource", DeleteIamResource)

	rg.GET("/types/:type/resources/:resource/roles/:role/users", ListIamResourceRole)
	rg.POST("/types/:type/resources/:resource/roles/:role/users", NewIamResourceRole)
	rg.DELETE("/types/:type/resources/:resource/roles/:role/users/:user", DeleteIamResourceRole)

	rg.GET("/types/:type/resources/:resource/actions/:action/users/:user", GetIamActionUser)
}

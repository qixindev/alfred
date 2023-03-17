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
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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

// ListIamResource godoc
//
//	@Summary		iam resource
//	@Schemes
//	@Description	get iam resource list
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
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

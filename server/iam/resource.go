package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/models/iam"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	get iam resource type list
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types [get]
func ListIamType(c *gin.Context) {
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}
	types, err := iam.ListResourceTypes(client.TenantId, client.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("list resource types err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, types)
}

// NewIamType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	new iam resource type
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types [post]
func NewIamType(c *gin.Context) {
	var typ models.ResourceType
	if err := c.BindJSON(&typ); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}
	t, err := iam.CreateResourceType(client.TenantId, client.Id, &typ)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("create resource type err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, t)
}

// DeleteIamType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	delete iam resource type
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type} [delete]
func DeleteIamType(c *gin.Context) {
	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}
	typeName := c.Param("type")
	var typ models.ResourceType
	if err = internal.TenantDB(c).First(&typ, "client_id = ? AND name = ?", client.Id, typeName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource type err: " + err.Error())
		return
	}
	if err = iam.DeleteResourceType(client.TenantId, typ.Id); err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("delete resource type err: " + err.Error())
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
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources [get]
func ListIamResource(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get type err: " + err.Error())
		return
	}
	resources, err := iam.ListResources(typ.TenantId, typ.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("list resource err: " + err.Error())
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
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources [post]
func NewIamResource(c *gin.Context) {
	var resource models.Resource
	if err := c.BindJSON(&resource); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get type err: " + err.Error())
		return
	}
	r, err := iam.CreateResource(typ.TenantId, typ.Id, &resource)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("create resource err: " + err.Error())
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
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource} [delete]
func DeleteIamResource(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get type err: " + err.Error())
		return
	}
	resourceName := c.Param("resource")
	var resource models.Resource
	if err = internal.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource err: " + err.Error())
		return
	}
	if err = iam.DeleteResource(typ.TenantId, resource.Id); err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("delete resource type err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

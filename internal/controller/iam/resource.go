package iam

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service/iam"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListIamType godoc
//
//	@Summary	iam resource type
//	@Schemes
//	@Description	get iam resource type list
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types [get]
func ListIamType(c *gin.Context) {
	client, err := GetClientFromCid(c)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get client err", true)
		return
	}

	types, err := iam.ListResourceTypes(client.TenantId, client.Id)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource type err", true)
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types [post]
func NewIamType(c *gin.Context) {
	var typ model.ResourceType
	if err := c.BindJSON(&typ); err != nil {
		resp.ErrorRequest(c, err, "bind new iam type err")
		return
	}
	client, err := GetClientFromCid(c)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}

	t, err := iam.CreateResourceType(client.TenantId, client.Id, typ)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "create resource type err")
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId} [delete]
func DeleteIamType(c *gin.Context) {
	clientId := c.Param("client")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	var typ model.ResourceType
	if err := internal.TenantDB(c).First(&typ, "client_id = ? AND id = ?", clientId, typeId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get resource type err")
		return
	}

	if err := iam.DeleteResourceType(tenant.Id, typ.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource type err")
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [get]
func ListIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resources, err := iam.ListResources(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource err", true)
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [post]
func NewIamResource(c *gin.Context) {
	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		resp.ErrorRequest(c, err, "bind new iam resource err")
		return
	}
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	typ, err := iam.GetIamType(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get resource type err")
		return
	}

	r, err := iam.CreateResource(tenant.Id, typ.Id, &resource)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "create resource err")
		return
	}
	c.JSON(http.StatusOK, r)
}

// UpdateIamResource godoc
//
//	@Summary		iam resource
//	@Schemes
//	@Description	update iam resource name
//	@Tags			iam-resource
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"typeId"
//	@Param			resourceId	path	string	true	"resourceId"
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId} [put]
func UpdateIamResource(c *gin.Context) {
	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		resp.ErrorRequest(c, err, "bind update iam resource err")
		return
	}

	tenant := internal.GetTenant(c)
	resource.Id = c.Param("resourceId")
	resource.TypeId = c.Param("typeId")
	r, err := iam.UpdateResource(tenant.Id, &resource)
	if err != nil {
		resp.ErrorSqlUpdate(c, err, "modify resource err")
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId} [delete]
func DeleteIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resourceId := c.Param("resourceId")
	resource, err := iam.GetIamResource(tenant.Id, typeId, resourceId)
	if err != nil {
		resp.ErrReqParaCustom(c, "no such resource")
		return
	}
	if err = iam.DeleteResource(tenant.Id, resource.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource err")
		return
	}
	c.Status(http.StatusNoContent)
}

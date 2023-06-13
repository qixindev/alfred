package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
	iam2 "accounts/server/service/iam"
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
		internal.ErrReqParaCustom(c, "no such client")
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}

	types, err := iam2.ListResourceTypes(client.TenantId, client.Id)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource type list")
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
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
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
		internal.ErrReqParaCustom(c, "no such client")
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}

	t, err := iam2.CreateResourceType(client.TenantId, client.Id, &typ)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to create resource type")
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
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId} [delete]
func DeleteIamType(c *gin.Context) {
	clientId := c.Param("client")
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	var typ models.ResourceType
	if err := internal.TenantDB(c).First(&typ, "client_id = ? AND id = ?", clientId, typeId).Error; err != nil {
		internal.ErrReqParaCustom(c, "no such resource type")
		global.LOG.Error("get resource type err: " + err.Error())
		return
	}

	if err := iam2.DeleteResourceType(tenant.Id, typ.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource type")
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
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [get]
func ListIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resources, err := iam2.ListResources(tenant.Id, typeId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource list")
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
//	@Param			typeId		path	string	true	"tenant"
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [post]
func NewIamResource(c *gin.Context) {
	var resource models.Resource
	if err := c.BindJSON(&resource); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	typ, err := iam2.GetIamType(tenant.Id, typeId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such resource type")
		global.LOG.Error("create resource err: " + err.Error())
		return
	}

	r, err := iam2.CreateResource(tenant.Id, typ.Id, &resource)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to create resource")
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
//	@Param			typeId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId} [delete]
func DeleteIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resourceId := c.Param("resourceId")
	resource, err := iam2.GetIamResource(tenant.Id, typeId, resourceId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such resource")
		return
	}
	if err = iam2.DeleteResource(tenant.Id, resource.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource")
		global.LOG.Error("delete resource err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

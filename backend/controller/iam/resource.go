package iam

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/service/iam"
	"github.com/gin-gonic/gin"
)

// ListIamType
// @Summary	获取资源类型列表
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types [get]
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
	resp.SuccessWithArrayData(c, types, 0)
}

// NewIamType
// @Summary	new iam resource type
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	iamBody		body	model.ResourceType	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types [post]
func NewIamType(c *gin.Context) {
	var typ model.ResourceType
	if err := c.BindJSON(&typ); err != nil {
		resp.ErrorRequest(c, err)
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
	resp.SuccessWithData(c, t)
}

// DeleteIamType
// @Summary	delete iam resource type
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId} [delete]
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
	resp.Success(c)
}

// ListIamResource
// @Summary	获取资源列表
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [get]
func ListIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resources, err := iam.ListResources(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource err", true)
		return
	}
	resp.SuccessWithArrayData(c, resources, 0)
}

// NewIamResource
// @Summary	new iam resource
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	iamBody		body	model.Resource	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources [post]
func NewIamResource(c *gin.Context) {
	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		resp.ErrorRequest(c, err)
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
	resp.SuccessWithData(c, r)
}

// UpdateIamResource
// @Summary	update iam resource name
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"typeId"
// @Param	resourceId	path	string	true	"resourceId"
// @Param	iamBody		body	model.Resource	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId} [put]
func UpdateIamResource(c *gin.Context) {
	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		resp.ErrorRequest(c, err)
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
	resp.SuccessWithData(c, r)
}

// DeleteIamResource
// @Summary	delete iam resource
// @Tags	iam-resource
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	resourceId	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId} [delete]
func DeleteIamResource(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	resourceId := c.Param("resourceId")
	resource, err := iam.GetIamResource(tenant.Id, typeId, resourceId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get resource err")
		return
	}
	if err = iam.DeleteResource(tenant.Id, resource.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource err")
		return
	}
	resp.Success(c)
}

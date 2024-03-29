package iam

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/req"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service/iam"
	"github.com/gin-gonic/gin"
)

// ListIamRole
// @Summary	获取角色列表
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"type id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles [get]
func ListIamRole(c *gin.Context) {
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	roles, err := iam.ListResourceTypeRoles(tenant.Id, typeId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource type role err", true)
		return
	}
	resp.SuccessWithArrayData(c, roles, 0)
}

// NewIamRole
// @Summary	new iam role
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	iamBody		body	model.ResourceTypeRole	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles [post]
func NewIamRole(c *gin.Context) {
	var role model.ResourceTypeRole
	if err := c.BindJSON(&role); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	r, err := iam.CreateResourceTypeRole(tenant.Id, typeId, &role)
	if err != nil {
		resp.ErrorSqlCreate(c, err, "create resource role err")
		return
	}
	resp.SuccessWithData(c, r)
}

// DeleteIamRole
// @Summary	delete iam role
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId} [delete]
func DeleteIamRole(c *gin.Context) {
	typeId := c.Param("typeId")
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	role, err := iam.GetIamRole(tenant.Id, typeId, roleId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get role err")
		return
	}

	if err = iam.DeleteResourceTypeRole(tenant.Id, role.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource role err")
		return
	}
	resp.Success(c)
}

// ListIamResourceRole
// @Summary	获取某个资源下某个角色的所有用户列表
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Param	resourceId	path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users [get]
func ListIamResourceRole(c *gin.Context) {
	resourceId := c.Param("resourceId")
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	roleUsers, err := iam.ListResourcesRoleUsers(tenant.Id, resourceId, roleId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resources role user err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(roleUsers, model.ResourceRoleUserDto), 0)
}

// NewIamResourceRole
// @Summary	new iam resource role
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Param	resourceId	path	string	true	"tenant"
// @Param	iamBody		body	[]model.ResourceRoleUser	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users [post]
func NewIamResourceRole(c *gin.Context) {
	var roleUser []model.ResourceRoleUser
	if err := c.BindJSON(&roleUser); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	roleId := c.Param("roleId")
	resourceId := c.Param("resourceId")
	tenant := internal.GetTenant(c)
	for i := 0; i < len(roleUser); i++ {
		if roleUser[i].ClientUserId == 0 {
			resp.ErrReqParaCustom(c, "client user id should not be empty")
			return
		}
		roleUser[i].RoleId = roleId
		roleUser[i].TenantId = tenant.Id
		roleUser[i].ResourceId = resourceId
	}

	if err := iam.CreateResourceRoleUser(tenant.Id, roleUser); err != nil {
		resp.ErrorSqlCreate(c, err, "create resource role user err")
		return
	}
	resp.SuccessWithData(c, roleUser)
}

// DeleteIamResourceRoleUser
// @Summary	delete iam resource role
// @Tags	iam-role
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	client		path	string	true	"client"	default(default)
// @Param	typeId		path	string	true	"tenant"
// @Param	resourceId	path	string	true	"tenant"
// @Param	roleId		path	string	true	"tenant"
// @Param	userId		path	string	true	"tenant"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users/{userId} [delete]
func DeleteIamResourceRoleUser(c *gin.Context) {
	resourceId := c.Param("resourceId")
	roleId := c.Param("roleId")
	sub := c.Param("userId")
	clientId := c.Param("client")
	tenant := internal.GetTenant(c)

	var roleUser model.ResourceRoleUser
	if err := global.DB.Table("resource_role_users as ru").Select("ru.id").
		Joins("LEFT JOIN client_users as cu ON ru.client_user_id = cu.id").
		Where("ru.tenant_id = ? AND cu.client_id = ? AND ru.resource_id = ? AND ru.role_id = ? AND cu.sub = ?",
			tenant.Id, clientId, resourceId, roleId, sub).
		First(&roleUser).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get resource role user")
		return
	}

	if err := iam.DeleteResourceRoleUser(tenant.Id, roleUser.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete resource role user err")
		return
	}
	resp.Success(c)
}

// CreateAllTypeRole
// @Summary	授权一类资源给用户
// @Tags	iam-role
// @Param	tenant		path	string		true	"tenant"	default(default)
// @Param	client		path	string		true	"client"	default(default)
// @Param	typeId		path	string		true	"type id"
// @Param	roleId		path	string		true	"role id"
// @Param	bd			body	[]integer	true	"client user id"
// @Success	200
// @Router	/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId}/auth [post]
func CreateAllTypeRole(c *gin.Context) {
	var in req.IamClientUser
	if err := internal.New(c).BindUri(&in).BindJson(&in.ClientUserId).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	tenant := internal.GetTenant(c)
	resources, err := iam.ListResources(tenant.Id, in.TypeId)
	if err != nil {
		resp.ErrorSqlSelect(c, err, "list resource err")
		return
	}

	roleUser := make([]model.ResourceRoleUser, 0)
	for _, cId := range in.ClientUserId {
		for _, resource := range resources {
			roleUser = append(roleUser, model.ResourceRoleUser{
				RoleId:       c.Param("roleId"),
				TenantId:     tenant.Id,
				ResourceId:   resource.Id,
				ClientUserId: cId.UserId,
			})
		}
	}

	if err = iam.CreateResourceRoleUser(tenant.Id, roleUser); err != nil {
		resp.ErrorSqlCreate(c, err, "create resource role user err")
		return
	}

	resp.Success(c)
}

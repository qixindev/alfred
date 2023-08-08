package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
	"accounts/server/service/iam"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// ListIamRole godoc
//
//	@Summary	iam role
//	@Schemes
//	@Description	get iam role list
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			typeId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles [get]
func ListIamRole(c *gin.Context) {
	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	roles, err := iam.ListResourceTypeRoles(tenant.Id, typeId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource type role list")
		global.LOG.Error("list resource type role err: " + err.Error())
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
//	@Param			typeId		path	string	true	"tenant"
//	@Param			iamBody		body	internal.IamNameRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles [post]
func NewIamRole(c *gin.Context) {
	var role models.ResourceTypeRole
	if err := c.BindJSON(&role); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	typeId := c.Param("typeId")
	tenant := internal.GetTenant(c)
	r, err := iam.CreateResourceTypeRole(tenant.Id, typeId, &role)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to create resource role")
		global.LOG.Error("create resource type role err: " + err.Error())
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
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/roles/{roleId} [delete]
func DeleteIamRole(c *gin.Context) {
	typeId := c.Param("typeId")
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	role, err := iam.GetIamRole(tenant.Id, typeId, roleId)
	if err != nil {
		internal.ErrReqParaCustom(c, "no such role")
		global.LOG.Error("get iam role err: ", zap.Error(err))
		return
	}

	if err = iam.DeleteResourceTypeRole(tenant.Id, role.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource role")
		global.LOG.Error("delete resource type role err: " + err.Error())
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
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users [get]
func ListIamResourceRole(c *gin.Context) {
	resourceId := c.Param("resourceId")
	roleId := c.Param("roleId")
	tenant := internal.GetTenant(c)
	roleUsers, err := iam.ListResourcesRoleUsers(tenant.Id, resourceId, roleId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resources role user list")
		global.LOG.Error("list resource role users err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(roleUsers, models.ResourceRoleUserDto))
}

// NewIamResourceRole godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	new iam resource role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Param			iamBody		body	[]internal.IamUserRequest	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users [post]
func NewIamResourceRole(c *gin.Context) {
	var roleUser []models.ResourceRoleUser
	if err := c.BindJSON(&roleUser); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	roleId := c.Param("roleId")
	resourceId := c.Param("resourceId")
	tenant := internal.GetTenant(c)
	for i := 0; i < len(roleUser); i++ {
		if roleUser[i].ClientUserId == 0 {
			internal.ErrReqParaCustom(c, "client user id should not be empty")
			return
		}
		roleUser[i].RoleId = roleId
		roleUser[i].TenantId = tenant.Id
		roleUser[i].ResourceId = resourceId
	}

	if err := iam.CreateResourceRoleUser(tenant.Id, roleUser); err != nil {
		internal.ErrorSqlResponse(c, "failed to create resource role user")
		global.LOG.Error("create resource role user err: " + err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// DeleteIamResourceRoleUser godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	delete iam resource role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			typeId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			userId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/roles/{roleId}/users/{userId} [delete]
func DeleteIamResourceRoleUser(c *gin.Context) {
	resourceId := c.Param("resourceId")
	roleId := c.Param("roleId")
	sub := c.Param("userId")
	clientId := c.Param("client")
	tenant := internal.GetTenant(c)

	var roleUser models.ResourceRoleUser
	if err := global.DB.Table("resource_role_users as ru").Select("ru.id").
		Joins("LEFT JOIN client_users as cu ON ru.client_user_id = cu.id").
		Where("ru.tenant_id = ? AND cu.client_id = ? AND ru.resource_id = ? AND ru.role_id = ? AND cu.sub = ?",
			tenant.Id, clientId, resourceId, roleId, sub).
		First(&roleUser).Error; err != nil {
		internal.ErrReqParaCustom(c, "no such resource role user")
		global.LOG.Error("delete resource role user err: " + err.Error())
		return
	}

	if err := iam.DeleteResourceRoleUser(tenant.Id, roleUser.Id); err != nil {
		internal.ErrorSqlResponse(c, "failed to delete resource role user")
		global.LOG.Error("delete resource role user err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateAllTypeRole godoc
//
//	@Summary		iam resource role
//	@Schemes
//	@Description	delete iam resource role
//	@Tags			iam-role
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			typeId		path	string	true	"tenant"
//	@Param			roleId		path	string	true	"tenant"
//	@Param			userId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/all/users/{userId}/roles/{roleId} [post]
func CreateAllTypeRole(c *gin.Context) {
	tenant := internal.GetTenant(c)
	typeId := c.Param("typeId")
	userId := c.Param("userId")
	userIdUint, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		internal.ErrReqParaCustom(c, "user id should be a number")
		global.LOG.Error("create all type role err: " + err.Error())
		return
	}

	resources, err := iam.ListResources(tenant.Id, typeId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get resource list")
		global.LOG.Error("list resource err: " + err.Error())
		return
	}
	for _, resource := range resources {
		roleUser := []models.ResourceRoleUser{{
			RoleId:       c.Param("roleId"),
			TenantId:     tenant.Id,
			ResourceId:   resource.Id,
			ClientUserId: uint(userIdUint),
		}}
		if err = iam.CreateResourceRoleUser(tenant.Id, roleUser); err != nil {
			internal.ErrorSqlResponse(c, "failed to create resource role user")
			global.LOG.Error("create resource role user err: " + err.Error())
			return
		}
	}

	internal.Success(c)
}

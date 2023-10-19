package iam

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service/iam"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
)

// IsUserActionPermission godoc
//
//	@Summary	用户是否拥有操作某个资源的权限
//	@Schemes
//	@Tags		iam-action
//	@Param		tenant		path	string	true	"tenant"	default(default)
//	@Param		client		path	string	true	"client"	default(default)
//	@Param		typeId		path	string	true	"tenant"
//	@Param		resourceId	path	string	true	"tenant"
//	@Param		actionId	path	string	true	"tenant"
//	@Param		user		path	string	true	"tenant"
//	@Success	200
//	@Router		/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/actions/{actionId}/users/{user} [get]
func IsUserActionPermission(c *gin.Context) {
	tenant := internal.GetTenant(c)
	resourceId := c.Param("resourceId")
	actionId := c.Param("actionId")
	userName := c.Param("user")
	clientId := c.Param("client")

	var clientUser model.ClientUser
	if err := internal.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", clientId, userName).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client user err")
		return
	}

	result, err := iam.CheckPermission(tenant.Id, clientUser.Id, resourceId, actionId)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "failed to check permission")
		return
	}

	resp.SuccessAuth(c, gin.H{"permission": result})
}

// GetIamActionResource godoc
//
//	@Summary	用户获取拥有某个操作的所有资源列表
//	@Schemes
//	@Tags		iam-action
//	@Param		tenant		path	string	true	"tenant"	default(default)
//	@Param		client		path	string	true	"client"	default(default)
//	@Param		typeId		path	string	true	"type id"
//	@Param		actionId	path	string	true	"action id"
//	@Param		user		path	string	true	"user sub id"
//	@Success	200
//	@Router		/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions/{actionId}/users/{user}/resources [get]
func GetIamActionResource(c *gin.Context) {
	typeId := c.Param("typeId")
	actionId := c.Param("actionId")
	user := c.Param("user")
	tenant := internal.GetTenant(c)
	res := make([]model.ResourceRoleUser, 0)

	if err := global.DB.Table("resource_role_users as rru").
		Select("rru.id", "rru.resource_id", "r.name resource_name", "rru.role_id", "rr.name role_name", "rru.client_user_id", "cu.sub").
		Joins("LEFT JOIN resources r ON r.id = rru.resource_id").
		Joins("LEFT JOIN resource_type_roles rr ON rr.id = rru.role_id").
		Joins("LEFT JOIN client_users cu ON cu.id = rru.client_user_id").
		Joins("LEFT JOIN resource_type_role_actions rtra ON rtra.role_id = rr.id").
		Where("rru.tenant_id = ? AND cu.sub = ? AND r.type_id = ? AND rtra.action_id = ?",
			tenant.Id, user, typeId, actionId).Find(&res).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list user's resources err", true)
		return
	}

	resp.SuccessWithArrayData(c, utils.Filter(res, model.ResourceRoleUserDto), 0)
}

// ListResourceUser godoc
//
//	@Summary	获取某个资源下所有用户列表
//	@Schemes
//	@Tags		iam-action
//	@Param		tenant		path	string	true	"tenant"	default(default)
//	@Param		client		path	string	true	"client"	default(default)
//	@Param		typeId		path	string	true	"tenant"
//	@Param		resourceId	path	string	true	"tenant"
//	@Success	200
//	@Router		/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/users [get]
func ListResourceUser(c *gin.Context) {
	typeId := c.Param("typeId")
	resourceId := c.Param("resourceId")
	tenant := internal.GetTenant(c)
	var res []model.ResourceRoleUser
	if err := global.DB.Table("resource_role_users as rru").
		Select("rru.id", "rru.role_id", "ro.name role_name", "rru.client_user_id", "cu.sub", "u.display_name").
		Joins("LEFT JOIN resources as r ON r.id = rru.resource_id").
		Joins("LEFT JOIN resource_type_roles ro ON ro.id = rru.role_id").
		Joins("LEFT JOIN client_users cu ON cu.id = rru.client_user_id").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("rru.tenant_id = ? AND rru.resource_id = ? AND r.type_id = ?", tenant.Id, resourceId, typeId).
		Find(&res).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list resource user err")
		return
	}

	resp.SuccessWithArrayData(c, utils.Filter(res, model.ResourceRoleUserDto), 0)
}

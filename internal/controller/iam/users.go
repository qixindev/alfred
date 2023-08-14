package iam

import (
	"accounts/internal/controller/internal"
	"accounts/internal/model"
	"accounts/internal/service/iam"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// IsUserActionPermission godoc
//
//	@Summary		iam action user
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			resourceId	path	string	true	"tenant"
//	@Param			actionId	path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/actions/{actionId}/users/{user} [get]
func IsUserActionPermission(c *gin.Context) {
	tenant := internal.GetTenant(c)
	resourceId := c.Param("resourceId")
	actionId := c.Param("actionId")
	userName := c.Param("user")
	clientId := c.Param("client")

	var clientUser model.ClientUser
	if err := internal.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", clientId, userName).Error; err != nil {
		internal.ErrReqParaCustom(c, "no such client user")
		global.LOG.Error("get client user err: " + err.Error())
		return
	}

	result, err := iam.CheckPermission(tenant.Id, clientUser.Id, resourceId, actionId)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to check permission")
		global.LOG.Error("check permission err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": result})
}

// GetIamActionResource godoc
//
//	@Summary		iam users roles
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			actionId	path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/actions/{actionId}/users/{user}/resources [get]
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
		internal.ErrorSqlResponse(c, "failed to get user's resources")
		global.LOG.Error("get resource err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.Filter(res, model.ResourceRoleUserDto))
}

// GetResourceUserList godoc
//
//	@Summary		iam users roles
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			client		path	string	true	"client"	default(default)
//	@Param			typeId		path	string	true	"tenant"
//	@Param			resourceId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{typeId}/resources/{resourceId}/users [get]
func GetResourceUserList(c *gin.Context) {
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
		internal.ErrorSqlResponse(c, "get resource user list err")
		global.LOG.Error("get resource err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.Filter(res, model.ResourceRoleUserDto))
}

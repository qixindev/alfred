package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/models/iam"
	"accounts/server/internal"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// IsUserActionPermission godoc
//
//	@Summary		iam action user
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			action		path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/resources/{resource}/actions/{action}/users/{user} [get]
func IsUserActionPermission(c *gin.Context) {
	resource, err := getResource(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource err: " + err.Error())
		return
	}
	action, err := getAction(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get action err: " + err.Error())
		return
	}

	client, err := GetClientFromCid(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client from cid err: " + err.Error())
		return
	}
	userName := c.Param("user")
	var clientUser models.ClientUser
	if err = internal.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", client.Id, userName).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get client user err: " + err.Error())
		return
	}

	result, err := iam.CheckPermission(resource.TenantId, clientUser.Id, resource.Id, action.Id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("check permission err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetIamActionResource godoc
//
//	@Summary		iam users roles
//	@Schemes
//	@Description	get iam action user
//	@Tags			iam-action
//	@Param			tenant		path	string	true	"tenant"
//	@Param			client		path	string	true	"tenant"
//	@Param			type		path	string	true	"tenant"
//	@Param			action		path	string	true	"tenant"
//	@Param			user		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/{type}/actions/{action}/users/{user}/resources [get]
func GetIamActionResource(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		internal.ErrorSqlResponse(c, "failed to get type")
		global.LOG.Error("get type err: " + err.Error())
		return
	}

	tenant := internal.GetTenant(c)
	res := make([]models.ResourceRoleUser, 0)
	if err = global.DB.Table("resource_role_users as rru").
		Select("r.name resource_name, rr.name role_name, cu.sub").
		Joins("LEFT JOIN resource_type_roles rr ON rr.id = rru.role_id").
		Joins("LEFT JOIN resource_type_role_actions rtra ON rtra.role_id = rr.id").
		Joins("LEFT JOIN resource_type_actions a ON a.id = rtra.action_id").
		Joins("LEFT JOIN resources r ON r.id = rru.resource_id").
		Joins("LEFT JOIN client_users cu ON cu.id = rru.client_user_id").
		Where("rru.tenant_id = ? AND cu.sub = ? AND r.type_id = ? AND a.name = ?",
			tenant.Id, c.Param("user"), typ.Id, c.Param("action")).Find(&res).Error; err != nil {
		internal.ErrorSqlResponse(c, "failed to get user's resources")
		global.LOG.Error("get resource err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.Filter(res, models.ResourceRoleUserDto))
}

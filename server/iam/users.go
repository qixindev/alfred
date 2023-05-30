package iam

import (
	"accounts/global"
	"accounts/models"
	"accounts/models/iam"
	"accounts/server/internal"
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

// GetIamUserRoles godoc
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
//	@Router			/accounts/{tenant}/iam/clients/{client}/types/:type/users/:user/roles [get]
func GetIamUserRoles(c *gin.Context) {
	typ, err := getType(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get type err: " + err.Error())
		return
	}

	tenant := internal.GetTenant(c)
	res := make([]models.ResourceRoleUser, 0)
	if err = global.DB.Table("resource_role_users as ru").
		Select("ru.id, ru.resource_id, ru.role_id, ru.client_user_id, ru.tenant_id").
		Joins("LEFT JOIN resources r ON ru.resource_id = r.id").
		Where("ru.tenant_id = ? AND ru.client_user_id = ? AND r.type_id = ?",
			tenant.Id, c.Param("user"), typ.Id).Find(&res).Error; err != nil {
		c.Status(http.StatusBadRequest)
		global.LOG.Error("get resource err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

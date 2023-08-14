package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/dto"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListGroups godoc
//
//	@Summary	group
//	@Schemes
//	@Description	list groups
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups [get]
func ListGroups(c *gin.Context) {
	var groups []model.Group
	if err := internal.TenantDB(c).Find(&groups).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(groups, model.Group2Dto))
}

// GetGroup godoc
//
//	@Summary	group
//	@Schemes
//	@Description	get groups
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups/{groupId} [get]
func GetGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, group.Dto())
}

// NewGroup godoc
//
//	@Summary	group
//	@Schemes
//	@Description	new groups
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups [post]
func NewGroup(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var group model.Group
	if err := c.BindJSON(&group); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	group.TenantId = tenant.Id
	if err := global.DB.Create(&group).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new group err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, group.Dto())
}

// UpdateGroup godoc
//
//	@Summary	group
//	@Schemes
//	@Description	update groups
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups/{groupId} [put]
func UpdateGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	var g model.Group
	if err := c.BindJSON(&g); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	group.Name = g.Name
	group.ParentId = g.ParentId
	if err := global.DB.Save(&group).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update group err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, group.Dto())
}

// DeleteGroup godoc
//
//	@Summary	group
//	@Schemes
//	@Description	delete groups
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups/{groupId} [delete]
func DeleteGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	if err := global.DB.Delete(&group).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete group err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// GetGroupMembers godoc
//
//	@Summary	group
//	@Schemes
//	@Description	get groups members
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups/{groupId}/member [get]
func GetGroupMembers(c *gin.Context) {
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}

	var members []dto.GroupMemberDto
	var groups []model.Group
	if err := internal.TenantDB(c).Find(&groups, "parent_id = ?", group.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	for _, g := range groups {
		members = append(members, g.GroupMemberDto())
	}

	var groupUsers []model.GroupUser
	if err := global.DB.Joins("User", "group_users.user_id = users.id AND group_users.tenant_id = users.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get group member err: " + err.Error())
		return
	}
	for _, u := range groupUsers {
		members = append(members, u.GroupMemberDto())
	}

	var groupDevices []model.GroupDevice
	if err := global.DB.Joins("Device", "group_devices.device_id = devices.id AND group_devices.tenant_id = devices.tenant_id").
		Find(&groupDevices, "group_devices.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get group device err: " + err.Error())
		return
	}
	for _, d := range groupDevices {
		members = append(members, d.GroupMemberDto())
	}

	c.JSON(http.StatusOK, members)
}

func AddAdminGroupsRoutes(rg *gin.RouterGroup) {
	rg.GET("/groups", ListGroups)
	rg.GET("/groups/:groupId", GetGroup)
	rg.POST("/groups", NewGroup)
	rg.PUT("/groups/:groupId", UpdateGroup)
	rg.DELETE("/groups/:groupId", DeleteGroup)
	rg.GET("/groups/:groupId/members", GetGroupMembers)
}

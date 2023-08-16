package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/dto"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
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
		resp.ErrorSqlSelect(c, err, "list groups err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(groups, model.Group2Dto), 0)
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
		resp.ErrorSqlFirst(c, err, "get group err")
		return
	}
	resp.SuccessWithData(c, group.Dto())
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
		resp.ErrorRequest(c, err, "bind new group err")
		return
	}
	group.TenantId = tenant.Id
	if err := global.DB.Create(&group).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "create group err")
		return
	}
	resp.SuccessWithData(c, group.Dto())
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
		resp.ErrorSqlFirst(c, err, "get group err")
		return
	}
	var g model.Group
	if err := c.BindJSON(&g); err != nil {
		resp.ErrorRequest(c, err, "bind update group err")
		return
	}
	group.Name = g.Name
	group.ParentId = g.ParentId
	if err := global.DB.Save(&group).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update group err")
		return
	}
	resp.SuccessWithData(c, group.Dto())
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
		resp.ErrorSqlFirst(c, err, "get group err")
		return
	}
	if err := global.DB.Delete(&group).Error; err != nil {
		resp.ErrorSqlDelete(c, err, "delete group err")
		return
	}
	resp.Success(c)
}

// ListGroupMembers godoc
//
//	@Summary	group
//	@Schemes
//	@Description	get groups members
//	@Tags			group
//	@Param			tenant	path	string	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/groups/{groupId}/member [get]
func ListGroupMembers(c *gin.Context) {
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get group err", true)
		return
	}

	var members []dto.GroupMemberDto
	var groups []model.Group
	if err := internal.TenantDB(c).Find(&groups, "parent_id = ?", group.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list group err")
		return
	}
	for _, g := range groups {
		members = append(members, g.GroupMemberDto())
	}

	var groupUsers []model.GroupUser
	if err := global.DB.Joins("User", "group_users.user_id = users.id AND group_users.tenant_id = users.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "get group member err")
		return
	}
	for _, u := range groupUsers {
		members = append(members, u.GroupMemberDto())
	}

	var groupDevices []model.GroupDevice
	if err := global.DB.Joins("Device", "group_devices.device_id = devices.id AND group_devices.tenant_id = devices.tenant_id").
		Find(&groupDevices, "group_devices.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list group device err")
		return
	}
	for _, d := range groupDevices {
		members = append(members, d.GroupMemberDto())
	}

	resp.SuccessWithArrayData(c, members, 0)
}

func AddAdminGroupsRoutes(rg *gin.RouterGroup) {
	rg.GET("/groups", ListGroups)
	rg.GET("/groups/:groupId", GetGroup)
	rg.POST("/groups", NewGroup)
	rg.PUT("/groups/:groupId", UpdateGroup)
	rg.DELETE("/groups/:groupId", DeleteGroup)
	rg.GET("/groups/:groupId/members", ListGroupMembers)
}

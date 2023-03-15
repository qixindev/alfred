package admin

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"accounts/utils"
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
//	@Router			/admin/{tenant}/groups [get]
func ListGroups(c *gin.Context) {
	var groups []models.Group
	if middlewares.TenantDB(c).Find(&groups).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(groups, models.Group2Dto))
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
//	@Router			/admin/{tenant}/groups/{groupId} [get]
func GetGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group models.Group
	if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
		c.Status(http.StatusNotFound)
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
//	@Router			/admin/{tenant}/groups [post]
func NewGroup(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var group models.Group
	err := c.BindJSON(&group)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	group.TenantId = tenant.Id
	if data.DB.Create(&group).Error != nil {
		c.Status(http.StatusConflict)
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
//	@Router			/admin/{tenant}/groups/{groupId} [put]
func UpdateGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group models.Group
	if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var g models.Group
	err := c.BindJSON(&g)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	group.Name = g.Name
	group.ParentId = g.ParentId
	if data.DB.Save(&group).Error != nil {
		c.Status(http.StatusInternalServerError)
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
//	@Router			/admin/{tenant}/groups/{groupId} [delete]
func DeleteGroup(c *gin.Context) {
	groupId := c.Param("groupId")
	var group models.Group
	if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if data.DB.Delete(&group).Error != nil {
		c.Status(http.StatusInternalServerError)
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
//	@Router			/admin/{tenant}/groups/{groupId}/member [get]
func GetGroupMembers(c *gin.Context) {
	groupId := c.Param("groupId")
	var group models.Group
	if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var members []dto.GroupMemberDto

	var groups []models.Group
	if middlewares.TenantDB(c).Find(&groups, "parent_id = ?", group.Id).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for _, g := range groups {
		members = append(members, g.GroupMemberDto())
	}

	var groupUsers []models.GroupUser
	if data.DB.Joins("User", "group_users.user_id = users.id AND group_users.tenant_id = users.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for _, u := range groupUsers {
		members = append(members, u.GroupMemberDto())
	}

	var groupDevices []models.GroupDevice
	if data.DB.Joins("Device", "group_devices.device_id = devices.id AND group_devices.tenant_id = devices.tenant_id").
		Find(&groupDevices, "group_devices.tenant_id = ? AND group_id = ?", group.TenantId, group.Id).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	for _, d := range groupDevices {
		members = append(members, d.GroupMemberDto())
	}

	c.JSON(http.StatusOK, members)
}

func addAdminGroupsRoutes(rg *gin.RouterGroup) {
	rg.GET("/groups", ListGroups)
	rg.GET("/groups/:groupId", GetGroup)
	rg.POST("/groups", NewGroup)
	rg.PUT("/groups/:groupId", UpdateGroup)
	rg.DELETE("/groups/:groupId", DeleteGroup)
	rg.GET("/groups/:groupId/members", GetGroupMembers)
}

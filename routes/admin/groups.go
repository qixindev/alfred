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

func addAdminGroupsRoutes(rg *gin.RouterGroup) {
	rg.GET("/groups", func(c *gin.Context) {
		var groups []models.Group
		if middlewares.TenantDB(c).Find(&groups).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(groups, models.Group2Dto))
	})

	rg.GET("/groups/:groupId", func(c *gin.Context) {
		groupId := c.Param("groupId")
		var group models.Group
		if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, group.Dto())
	})

	rg.POST("/groups", func(c *gin.Context) {
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
	})

	rg.PUT("/groups/:groupId", func(c *gin.Context) {
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
	})

	rg.DELETE("/groups/:groupId", func(c *gin.Context) {
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
	})

	rg.GET("/groups/:groupId/members", func(c *gin.Context) {
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
	})
}

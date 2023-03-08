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

func addAdminTenantsRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", func(c *gin.Context) {
		var tenants []models.Tenant
		if data.DB.Find(&tenants).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(tenants, models.Tenant2Dto))
	})

	rg.GET("/tenants/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		var tenant models.Tenant
		if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, tenant.Dto())
	})

	rg.POST("/tenants", func(c *gin.Context) {
		var tenant models.Tenant
		err := c.BindJSON(&tenant)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if data.DB.Create(&tenant).Error != nil {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusOK, tenant.Dto())
	})

	rg.PUT("/tenants/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		var tenant models.Tenant
		if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var t models.Tenant
		err := c.BindJSON(&t)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		tenant.Name = t.Name
		if data.DB.Save(&tenant).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, tenant.Dto())
	})

	rg.DELETE("/tenants/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		var tenant models.Tenant
		if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if data.DB.Delete(&tenant).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/admin/:tenant/devices/:deviceId/groups", func(c *gin.Context) {
		deviceId := c.Param("deviceId")
		var device models.Device
		if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var groupDevices []models.GroupDevice
		if data.DB.Joins("Group", "group_devices.group_id = groups.id AND group_devices.tenant_id = groups.tenant_id").
			Find(&groupDevices, "group_devices.tenant_id = ? AND device_id = ?", device.TenantId, device.Id).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		groups := utils.Filter(groupDevices, func(gd models.GroupDevice) dto.GroupMemberDto {
			return dto.GroupMemberDto{
				Id:   gd.GroupId,
				Name: gd.Group.Name,
			}
		})
		c.JSON(http.StatusOK, groups)
	})

	rg.PUT("/devices/:deviceId/groups/:groupId", func(c *gin.Context) {
		deviceId := c.Param("deviceId")
		var device models.Device
		if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		groupId := c.Param("groupId")
		var group models.Group
		if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var groupDevice models.GroupDevice
		if middlewares.TenantDB(c).First(groupDevice, "group_id = ? AND device_id = ?", group.Id, device.Id).Error != nil {
			// Not found, create one.
			groupDevice.DeviceId = device.Id
			groupDevice.GroupId = group.Id
			groupDevice.TenantId = device.TenantId
		} else {
			// Found, update it.
			if middlewares.TenantDB(c).Save(&groupDevice).Error != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		c.JSON(http.StatusOK, groupDevice.GroupMemberDto())
	})

	rg.DELETE("/admin/:tenant/devices/:deviceId/groups/:groupId", func(c *gin.Context) {
		deviceId := c.Param("deviceId")
		var device models.Device
		if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		groupId := c.Param("groupId")
		var groupDevice models.GroupDevice
		if middlewares.TenantDB(c).First(&groupDevice, "device_id = ? AND group_id = ?", device.Id, groupId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if middlewares.TenantDB(c).Delete(&groupDevice).Error != nil {
			c.Status(http.StatusInternalServerError)
		}
		c.Status(http.StatusNoContent)
	})
}

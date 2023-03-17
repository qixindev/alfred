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

// ListDevices godoc
//
//	@Summary	device
//	@Schemes
//	@Description	list device
//	@Tags			device
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices [get]
func ListDevices(c *gin.Context) {
	var devices []models.Device
	if middlewares.TenantDB(c).Find(&devices).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(devices, models.Device2Dto))
}

// GetDevice godoc
//
//	@Summary	device
//	@Schemes
//	@Description	get device
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId} [get]
func GetDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device models.Device
	if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, device.Dto())
}

// NewDevice godoc
//
//	@Summary	device
//	@Schemes
//	@Description	new device
//	@Tags			device
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices [post]
func NewDevice(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var device models.Device
	err := c.BindJSON(&device)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	device.TenantId = tenant.Id
	if data.DB.Create(&device).Error != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, device.Dto())
}

// UpdateDevice godoc
//
//	@Summary	device
//	@Schemes
//	@Description	update device
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId} [put]
func UpdateDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device models.Device
	if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var d models.Device
	err := c.BindJSON(&d)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	device.Name = d.Name
	if data.DB.Save(&device).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, device.Dto())
}

// DeleteDevice godoc
//
//	@Summary	device
//	@Schemes
//	@Description	delete device
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId} [delete]
func DeleteDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device models.Device
	if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if data.DB.Delete(&device).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetDeviceGroups godoc
//
//	@Summary	device groups
//	@Schemes
//	@Description	list device groups
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId}/groups [get]
func GetDeviceGroups(c *gin.Context) {
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
}

// NewDeviceGroup godoc
//
//	@Summary	device groups
//	@Schemes
//	@Description	new device groups
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId}/groups [post]
func NewDeviceGroup(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var deviceGroup models.GroupDevice
	if err := c.BindJSON(&deviceGroup); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var device models.Device
	if middlewares.TenantDB(c).First(&device, "id = ?", deviceId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	deviceGroup.TenantId = device.TenantId
	deviceGroup.DeviceId = device.Id
	if data.DB.Create(&deviceGroup).Error != nil {
		c.Status(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, deviceGroup.Dto())
}

// UpdateDeviceGroup godoc
//
//	@Summary	device groups
//	@Schemes
//	@Description	update device groups
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Param			groupId		path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId}/groups/{groupId} [get]
func UpdateDeviceGroup(c *gin.Context) {
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
}

// DeleteDeviceGroup godoc
//
//	@Summary	device groups
//	@Schemes
//	@Description	delete device groups
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Param			groupId		path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/devices/{deviceId}/groups/{groupId} [delete]
func DeleteDeviceGroup(c *gin.Context) {
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
}

func addAdminDevicesRoutes(rg *gin.RouterGroup) {
	rg.GET("/devices", ListDevices)
	rg.GET("/devices/:deviceId", GetDevice)
	rg.POST("/devices", NewDevice)
	rg.PUT("/devices/:deviceId", UpdateDevice)
	rg.DELETE("/devices/:deviceId", DeleteDevice)
	rg.GET("/devices/:deviceId/groups", GetDeviceGroups)
	rg.POST("/devices/:deviceId/groups", NewDeviceGroup)
	rg.PUT("/devices/:deviceId/groups/:groupId", UpdateDeviceGroup)
	rg.DELETE("/devices/:deviceId/groups/:groupId", DeleteDeviceGroup)
}

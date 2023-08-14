package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/dto"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// ListDevices godoc
//
//	@Summary	device
//	@Schemes
//	@Description	list device
//	@Tags			device
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/devices [get]
func ListDevices(c *gin.Context) {
	var devices []model.Device
	if err := internal.TenantDB(c).Find(&devices).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get device err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(devices, model.Device2Dto))
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId} [get]
func GetDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/devices [post]
func NewDevice(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var device model.Device
	if err := c.BindJSON(&device); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	if device.Id == "" {
		device.Id = uuid.NewString()
	}
	device.TenantId = tenant.Id
	if err := global.DB.Create(&device).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new device err: " + err.Error())
		return
	}

	secret := model.DeviceSecret{
		DeviceId: device.Id,
		Name:     "default",
		Secret:   uuid.NewString(),
		TenantId: tenant.Id,
	}
	if err := internal.TenantDB(c).Create(&secret).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new device secret err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"id":                 device.Id,
		"device_name":        device.Name,
		"device_secret":      secret.Secret,
		"device_secret_name": secret.Name,
	})
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId} [put]
func UpdateDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}
	var d model.Device
	if err := c.BindJSON(&d); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	device.Name = d.Name
	if err := global.DB.Save(&device).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update device err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId} [delete]
func DeleteDevice(c *gin.Context) {
	deviceId := c.Param("deviceId")
	tenant := internal.GetTenant(c)
	if err := service.DeleteDevice(tenant.Id, deviceId); err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete device err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// ListDeviceSecret godoc
//
//	@Summary	get client secrets
//	@Schemes
//	@Description	get client secrets
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/secrets [get]
func ListDeviceSecret(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}
	var secrets []model.DeviceSecret
	if err := internal.TenantDB(c).Find(&secrets, "device_id = ?", device.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get device secret err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(secrets, model.DeviceSecret2Dto))
}

// NewDeviceSecret godoc
//
//	@Summary	get client secrets
//	@Schemes
//	@Description	get client secrets
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/secrets [post]
func NewDeviceSecret(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}
	var secret model.DeviceSecret
	if err := c.BindJSON(&secret); err != nil {
		resp.ErrReqPara(c, err)
		return
	}
	secret.DeviceId = device.Id
	secret.TenantId = device.TenantId
	if err := global.DB.Create(&secret).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("create device secret err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, secret.Dto())
}

// DeleteDeviceSecret godoc
//
//	@Summary	get client secrets
//	@Schemes
//	@Description	get client secrets
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			secretId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/secrets/{secretId} [delete]
func DeleteDeviceSecret(c *gin.Context) {
	deviceId := c.Param("deviceId")
	secretId := c.Param("secretId")
	tenant := internal.GetTenant(c)
	var secret model.DeviceSecret
	if err := internal.TenantDB(c).First(&secret, "tenant_id = ? AND device_id = ? AND id = ?", tenant.Id, deviceId, secretId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device secret err: " + err.Error())
		return
	}

	if err := global.DB.Delete(&secret).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete client secret err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/groups [get]
func GetDeviceGroups(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}
	var groupDevices []model.GroupDevice
	if err := global.DB.Joins("Group", "group_devices.group_id = groups.id AND group_devices.tenant_id = groups.tenant_id").
		Find(&groupDevices, "group_devices.tenant_id = ? AND device_id = ?", device.TenantId, device.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get device group err: " + err.Error())
		return
	}
	groups := utils.Filter(groupDevices, func(gd model.GroupDevice) dto.GroupMemberDto {
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/groups [post]
func NewDeviceGroup(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var deviceGroup model.GroupDevice
	if err := c.BindJSON(&deviceGroup); err != nil {
		resp.ErrReqPara(c, err)
		return
	}

	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}

	deviceGroup.TenantId = device.TenantId
	deviceGroup.DeviceId = device.Id
	if err := global.DB.Create(&deviceGroup).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new device group err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/groups/{groupId} [put]
func UpdateDeviceGroup(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}

	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}

	var groupDevice model.GroupDevice
	if err := internal.TenantDB(c).First(groupDevice, "group_id = ? AND device_id = ?", group.Id, device.Id).Error; err != nil {
		global.LOG.Error("get group device err: " + err.Error())
		// Not found, create one.
		groupDevice.DeviceId = device.Id
		groupDevice.GroupId = group.Id
		groupDevice.TenantId = device.TenantId
	} else {
		// Found, update it.
		if err := internal.TenantDB(c).Save(&groupDevice).Error; err != nil {
			c.Status(http.StatusInternalServerError)
			global.LOG.Error("get device group err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/devices/{deviceId}/groups/{groupId} [delete]
func DeleteDeviceGroup(c *gin.Context) {
	deviceId := c.Param("deviceId")
	var device model.Device
	if err := internal.TenantDB(c).First(&device, "id = ?", deviceId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get device err: " + err.Error())
		return
	}

	groupId := c.Param("groupId")
	var groupDevice model.GroupDevice
	if err := internal.TenantDB(c).First(&groupDevice, "device_id = ? AND group_id = ?", device.Id, groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group device err: " + err.Error())
		return
	}

	if err := internal.TenantDB(c).Delete(&groupDevice).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete device group err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// VerifyDeviceCode godoc
//
//	@Summary	device code
//	@Schemes
//	@Description	delete device groups
//	@Tags			device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			userCode	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/devices/code/{userCode} [post]
func VerifyDeviceCode(c *gin.Context) {
	userCode := c.Param("userCode")
	deviceCode := model.DeviceCode{}
	if err := internal.TenantDB(c).Where("user_code = ?", userCode).First(&deviceCode).Error; err != nil {
		c.String(http.StatusInternalServerError, "failed to get user code")
		global.LOG.Error("set device code err: " + err.Error())
		return
	}

	if deviceCode.CreatedAt.Add(2 * time.Minute).Before(time.Now()) {
		c.String(http.StatusGone, "user code expired")
		service.ClearDeviceCode(userCode)
		return
	}
	if err := internal.TenantDB(c).Table("device_codes").Where("user_code = ?", userCode).Update("status", "verified").Error; err != nil {
		c.String(http.StatusInternalServerError, "failed to verify user code")
		global.LOG.Error("set device code err: " + err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func AddAdminDevicesRoutes(rg *gin.RouterGroup) {
	rg.GET("/devices", ListDevices)
	rg.GET("/devices/:deviceId", GetDevice)
	rg.POST("/devices", NewDevice)
	rg.PUT("/devices/:deviceId", UpdateDevice)
	rg.DELETE("/devices/:deviceId", DeleteDevice)

	rg.GET("/devices/:deviceId/secrets", ListDeviceSecret)
	rg.POST("/devices/:deviceId/secrets", NewDeviceSecret)
	rg.DELETE("/devices/:deviceId/secret/:secretId", DeleteDeviceSecret)

	rg.GET("/devices/:deviceId/groups", GetDeviceGroups)
	rg.POST("/devices/:deviceId/groups", NewDeviceGroup)
	rg.PUT("/devices/:deviceId/groups/:groupId", UpdateDeviceGroup)
	rg.DELETE("/devices/:deviceId/groups/:groupId", DeleteDeviceGroup)

	rg.POST("/devices/code/:userCode", VerifyDeviceCode)
}

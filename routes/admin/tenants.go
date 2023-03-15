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

// ListTenants godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	list tenants
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/admin/tenants [get]
func ListTenants(c *gin.Context) {
	var tenants []models.Tenant
	if data.DB.Find(&tenants).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(tenants, models.Tenant2Dto))
}

// GetTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	get tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/admin/tenants/{tenantId} [get]
func GetTenant(c *gin.Context) {
	tenantId := c.Param("tenantId")
	var tenant models.Tenant
	if data.DB.First(&tenant, "id = ?", tenantId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, tenant.Dto())
}

// NewTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	new tenants
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/admin/tenants [post]
func NewTenant(c *gin.Context) {
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
}

// UpdateTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	update tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/admin/tenants/{tenantId} [put]
func UpdateTenant(c *gin.Context) {
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
}

// DeleteTenant godoc
//
//	@Summary	tenants
//	@Schemes
//	@Description	delete tenants
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"
//	@Param			tenantId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenant}/admin/tenants/{tenantId} [delete]
func DeleteTenant(c *gin.Context) {
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
}

// GetAdminDeviceGroup godoc
//
//	@Summary	admin device
//	@Schemes
//	@Description	get admin device group
//	@Tags			admin-device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenants}/devices/{deviceId}/groups [get]
func GetAdminDeviceGroup(c *gin.Context) {
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

// UpdateAdminDeviceGroups godoc
//
//	@Summary	admin device
//	@Schemes
//	@Description	update admin device group
//	@Tags			admin-device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenants}/devices/{deviceId}/groups [put]
func UpdateAdminDeviceGroups(c *gin.Context) {
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

// DeleteAdminDeviceGroup godoc
//
//	@Summary	admin device
//	@Schemes
//	@Description	delete admin device group
//	@Tags			admin-device
//	@Param			tenant		path	string	true	"tenant"
//	@Param			deviceId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/admin/{tenants}/devices/{deviceId}/groups [delete]
func DeleteAdminDeviceGroup(c *gin.Context) {
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

func addAdminTenantsRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenants", ListTenants)
	rg.GET("/tenants/:tenantId", GetTenant)
	rg.POST("/tenants", NewTenant)
	rg.PUT("/tenants/:tenantId", UpdateTenant)
	rg.DELETE("/tenants/:tenantId", DeleteTenant)

	rg.GET("/tenants/devices/:deviceId/groups", GetAdminDeviceGroup)
	rg.PUT("/tenants/devices/:deviceId/groups/:groupId", UpdateAdminDeviceGroups)
	rg.DELETE("/tenants/devices/:deviceId/groups/:groupId", DeleteAdminDeviceGroup)
}

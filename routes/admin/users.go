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

// ListUsers godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user list
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [get]
func ListUsers(c *gin.Context) {
	var users []models.User
	if middlewares.TenantDB(c).Find(&users).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(users, models.User2AdminDto))
}

// GetUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [get]
func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, user.AdminDto())
}

// NewUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	new user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [post]
func NewUser(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user.TenantId = tenant.Id
	if data.DB.Create(&user).Error != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, user.AdminDto())
}

// UpdateUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var u models.User
	err := c.BindJSON(&u)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user.Username = u.Username
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.DisplayName = u.DisplayName
	user.Email = u.Email
	user.EmailVerified = u.EmailVerified
	user.Phone = u.Phone
	user.PhoneVerified = u.PhoneVerified
	user.TwoFactorEnabled = u.TwoFactorEnabled
	user.Disabled = u.Disabled
	if data.DB.Save(&user).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, user.AdminDto())
}

// DeleteUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	delete user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [delete]
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if data.DB.Delete(&user).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetUserGroups godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups [get]
func GetUserGroups(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var groupUsers []models.GroupUser
	if data.DB.Joins("Group", "group_users.group_id = groups.id AND group_users.tenant_id = groups.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND user_id = ?", user.TenantId, user.Id).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	groups := utils.Filter(groupUsers, func(gu models.GroupUser) dto.GroupMemberDto {
		return dto.GroupMemberDto{
			Id:   gu.GroupId,
			Name: gu.Group.Name,
			Role: gu.Role,
		}
	})
	c.JSON(http.StatusOK, groups)
}

// NewUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups [post]
func NewUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var groupUser models.GroupUser
	if err := c.BindJSON(&groupUser); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	groupUser.TenantId = user.TenantId
	groupUser.UserId = user.Id
	groupUser.Role = user.Role
	if data.DB.Create(&groupUser).Error != nil {
		c.Status(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, groupUser.Dto())
}

// UpdateUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [get]
func UpdateUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	groupId := c.Param("groupId")
	var group models.Group
	if middlewares.TenantDB(c).First(&group, "id = ?", groupId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var gu dto.GroupMemberDto
	if c.BindJSON(&gu) != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var groupUser models.GroupUser
	if middlewares.TenantDB(c).First(groupUser, "group_id = ? AND user_id = ?", group.Id, user.Id).Error != nil {
		// Not found, create one.
		groupUser.UserId = user.Id
		groupUser.GroupId = group.Id
		groupUser.TenantId = user.TenantId
		groupUser.Role = gu.Role
	} else {
		// Found, update it.
		groupUser.Role = gu.Role
		if middlewares.TenantDB(c).Save(&groupUser).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusOK, groupUser.GroupMemberDto())
}

// DeleteUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"
//	@Param			userId	path	integer	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [delete]
func DeleteUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var user models.User
	if middlewares.TenantDB(c).First(&user, "id = ?", userId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	groupId := c.Param("groupId")
	var groupUser models.GroupUser
	if middlewares.TenantDB(c).First(&groupUser, "user_id = ? AND group_id = ?", user.Id, groupId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if middlewares.TenantDB(c).Delete(&groupUser).Error != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusNoContent)
}

func addAdminUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", ListUsers)
	rg.GET("/users/:userId", GetUser)
	rg.POST("/users", NewUser)
	rg.PUT("/users/:userId", UpdateUser)
	rg.DELETE("/users/:userId", DeleteUser)
	rg.GET("/users/:userId/groups", GetUserGroups)
	rg.POST("/users/:userId/groups", NewUserGroup)
	rg.PUT("/users/:userId/groups/:groupId", UpdateUserGroup)
	rg.DELETE("/users/:userId/groups/:groupId", DeleteUserGroup)
}

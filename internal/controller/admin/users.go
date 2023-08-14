package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/dto"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListUsers godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user list
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [get]
func ListUsers(c *gin.Context) {
	var users []model.User
	if err := internal.TenantDB(c).Find(&users).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get tenant users err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(users, model.User2AdminDto))
}

// GetUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [get]
func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
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
	tenant := internal.GetTenant(c)
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	if user.PasswordHash == "" {
		c.String(http.StatusBadRequest, "password should not be null")
		return
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("hashPassword err: " + err.Error())
		return
	}

	user.TenantId = tenant.Id
	user.PasswordHash = hash
	if err = global.DB.Create(&user).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new tenant user err: " + err.Error())
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
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}
	var u model.User
	if err := c.BindJSON(&u); err != nil {
		internal.ErrReqPara(c, err)
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
	user.Role = u.Role
	if err := global.DB.Save(&user).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update tenant user err: " + err.Error())
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
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}

	if err := service.DeleteUser(user.Id); err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete tenant user err: " + err.Error())
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
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}
	var groupUsers []model.GroupUser
	if err := global.DB.Joins("Group", "group_users.group_id = groups.id AND group_users.tenant_id = groups.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND user_id = ?", user.TenantId, user.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get tenant user groups err: " + err.Error())
		return
	}
	groups := utils.Filter(groupUsers, func(gu model.GroupUser) dto.GroupMemberDto {
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
	var groupUser model.GroupUser
	if err := c.BindJSON(&groupUser); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}

	groupUser.TenantId = user.TenantId
	groupUser.UserId = user.Id
	groupUser.Role = user.Role
	if err := global.DB.Create(&groupUser).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("create tenant user group err: " + err.Error())
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
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group err: " + err.Error())
		return
	}
	var gu dto.GroupMemberDto
	if err := c.BindJSON(&gu); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	var groupUser model.GroupUser
	if err := internal.TenantDB(c).First(groupUser, "group_id = ? AND user_id = ?", group.Id, user.Id).Error; err != nil {
		global.LOG.Error("get group user err: " + err.Error())
		// Not found, create one.
		groupUser.UserId = user.Id
		groupUser.GroupId = group.Id
		groupUser.TenantId = user.TenantId
		groupUser.Role = gu.Role
	} else {
		// Found, update it.
		groupUser.Role = gu.Role
		if err := internal.TenantDB(c).Save(&groupUser).Error; err != nil {
			c.Status(http.StatusInternalServerError)
			global.LOG.Error("get tenant user group err: " + err.Error())
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
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get user err: " + err.Error())
		return
	}
	groupId := c.Param("groupId")
	var groupUser model.GroupUser
	if err := internal.TenantDB(c).First(&groupUser, "user_id = ? AND group_id = ?", user.Id, groupId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get group user err: " + err.Error())
		return
	}
	if err := internal.TenantDB(c).Delete(&groupUser).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete user group err: " + err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func AddAdminUsersRoutes(rg *gin.RouterGroup) {
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

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
		resp.ErrorSqlSelect(c, err, "list tenant users err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(users, model.User2AdminDto), 0)
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
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// NewUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	new user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			bd		body	object	true	"user body"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [post]
func NewUser(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		resp.ErrorRequestWithMsg(c, err, "bind new user err")
		return
	}
	if user.PasswordHash == "" {
		resp.ErrorRequestWithMsg(c, nil, "password should not be null")
		return
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}

	user.TenantId = tenant.Id
	user.PasswordHash = hash
	if err = global.DB.Create(&user).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "new tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// UpdateUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	var u model.User
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequestWithMsg(c, err, "bind update user err")
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
		resp.ErrorSqlUpdate(c, err, "update tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// DeleteUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	delete user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [delete]
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}

	if err := service.DeleteUser(user.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete tenant user err")
		return
	}
	c.Status(http.StatusNoContent)
}

// ListUserGroups godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups [get]
func ListUserGroups(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err", true)
		return
	}
	var groupUsers []model.GroupUser
	if err := global.DB.Joins("Group", "group_users.group_id = groups.id AND group_users.tenant_id = groups.tenant_id").
		Find(&groupUsers, "group_users.tenant_id = ? AND user_id = ?", user.TenantId, user.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list tenant user groups err", true)
		return
	}
	groups := utils.Filter(groupUsers, func(gu model.GroupUser) dto.GroupMemberDto {
		return dto.GroupMemberDto{
			Id:   gu.GroupId,
			Name: gu.Group.Name,
			Role: gu.Role,
		}
	})
	resp.SuccessWithArrayData(c, groups, 0)
}

// NewUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups [post]
func NewUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var groupUser model.GroupUser
	if err := c.BindJSON(&groupUser); err != nil {
		resp.ErrorRequestWithMsg(c, err, "bind new user group err")
		return
	}

	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}

	groupUser.TenantId = user.TenantId
	groupUser.UserId = user.Id
	groupUser.Role = user.Role
	if err := global.DB.Create(&groupUser).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "create tenant group user err")
		return
	}

	resp.SuccessWithData(c, groupUser.Dto())
}

// UpdateUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [get]
func UpdateUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	groupId := c.Param("groupId")
	var group model.Group
	if err := internal.TenantDB(c).First(&group, "id = ?", groupId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get group err")
		return
	}
	var gu dto.GroupMemberDto
	if err := c.BindJSON(&gu); err != nil {
		resp.ErrorRequestWithMsg(c, err, "bind update user group err")
		return
	}
	var groupUser model.GroupUser
	if err := internal.TenantDB(c).First(groupUser, "group_id = ? AND user_id = ?", group.Id, user.Id).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get group user err")
		return
	}

	groupUser.Role = gu.Role
	if err := internal.TenantDB(c).Save(&groupUser).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "get tenant user group err")
		return
	}

	resp.SuccessWithData(c, groupUser.GroupMemberDto())
}

// DeleteUserGroup godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user groups
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Param			groupId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [delete]
func DeleteUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	groupId := c.Param("groupId")
	var groupUser model.GroupUser
	if err := internal.TenantDB(c).First(&groupUser, "user_id = ? AND group_id = ?", user.Id, groupId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get group user err")
		return
	}
	if err := internal.TenantDB(c).Delete(&groupUser).Error; err != nil {
		resp.ErrorSqlDelete(c, err, "delete group user err")
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

	rg.GET("/users/:userId/groups", ListUserGroups)
	rg.POST("/users/:userId/groups", NewUserGroup)
	rg.PUT("/users/:userId/groups/:groupId", UpdateUserGroup)
	rg.DELETE("/users/:userId/groups/:groupId", DeleteUserGroup)
}

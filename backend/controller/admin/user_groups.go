package admin

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/dto"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ListUserGroups
// @Summary	get user groups
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId}/groups [get]
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

// NewUserGroup
// @Summary	get user groups
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId}/groups [post]
func NewUserGroup(c *gin.Context) {
	userId := c.Param("userId")
	var groupUser model.GroupUser
	if err := c.BindJSON(&groupUser); err != nil {
		resp.ErrorRequest(c, err)
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

// UpdateUserGroup
// @Summary	update user groups
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Param	groupId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [get]
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
		resp.ErrorRequest(c, err)
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

// DeleteUserGroup
// @Summary	update user groups
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Param	groupId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId}/groups/{groupId} [delete]
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
	resp.Success(c)
}

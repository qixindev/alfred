package admin

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ListUsers
// @Summary	get user list
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200
// @Router	/accounts/admin/{tenant}/users [get]
func ListUsers(c *gin.Context) {
	var users []model.User
	if err := internal.TenantDB(c).Find(&users).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list tenant users err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(users, model.User2AdminDto), 0)
}

// GetUser
// @Summary	get user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId} [get]
func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// NewUser
// @Summary	new user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	bd	body	object	true	"user body"
// @Success	200
// @Router	/accounts/admin/{tenant}/users [post]
func NewUser(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if exist, err := service.IsUserPhoneOrEmailExist(user, tenant.Id); err != nil {
		resp.ErrorSqlFirst(c, err, "service.IsUserPhoneOrEmailExist err")
		return
	} else if exist {
		resp.ErrorConflict(c, err, "user phone or email is already exist.")
		return
	}

	if user.PasswordHash == "" {
		resp.ErrorRequestWithMsg(c, "password should not be null")
		return
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}

	user.TenantId = tenant.Id
	user.PasswordHash = hash
	user.From = "admin-create"
	user.Meta = "{}"
	if err = global.DB.Create(&user).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "new tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// UpdateUser
// @Summary	update user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"user id"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	var u model.User
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	hash, err := utils.HashPassword(u.PasswordHash)
	if err != nil {
		resp.ErrorUnauthorized(c, nil, "hashPassword err")
		return
	}

	user.Username = u.Username
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.DisplayName = u.DisplayName
	user.PasswordHash = hash
	user.Email = u.Email
	user.Phone = u.Phone
	user.Disabled = u.Disabled
	user.Role = u.Role
	user.Avatar = u.Avatar
	if err = global.DB.Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// DeleteUser
// @Summary	delete user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	userId	path	integer	true	"tenant"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/{userId} [delete]
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
	resp.Success(c)
}

// MarkUserFrom
// @Summary	mark user from
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	bd	body	object	true	"user body"
// @Success	200
// @Router	/accounts/admin/{tenant}/users/from [post]
func MarkUserFrom(c *gin.Context) {
	var request struct {
		IdList []int64 `json:"idList"`
		From   string  `json:"from"`
	}
	if err := c.BindJSON(&request); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if len(request.IdList) == 0 {
		resp.ErrorRequestWithMsg(c, "id list should not be null")
		return
	}
	if request.From == "" {
		resp.ErrorRequestWithMsg(c, "from should not be null")
		return
	}
	for _, id := range request.IdList {
		var user model.User
		if err := internal.TenantDB(c).First(&user, "id = ?", id).Error; err != nil {
			resp.ErrorSqlFirst(c, err, "get user err")
			return
		}
		user.From = request.From
		if err := global.DB.Save(&user).Error; err != nil {
			resp.ErrorSqlUpdate(c, err, "update tenant user err")
			return
		}
	}
	resp.Success(c)
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

	rg.POST("/users/from", MarkUserFrom)
}

package admin

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
	"io"
)

type ModifyPassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ListClientUsers
// @Summary	get client user list
// @Tags	client-user
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	clientId	path	string	true	"client id"	default(default)
// @Success	200
// @Router	/accounts/admin/{tenant}/clients/{clientId}/users [get]
func ListClientUsers(c *gin.Context) {
	var clientUser []struct {
		Sub      string `json:"sub"`
		ClientId string `json:"clientId"`
		model.User
	}
	clientId := c.Param("clientId")
	if err := global.DB.Table("client_users cu").
		Select("cu.id, cu.sub sub, cu.client_id, u.username username, u.phone, u.email, u.first_name, u.last_name, u.display_name, u.role, u.avatar, u.from").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ?", internal.GetTenant(c).Id, clientId).
		Find(&clientUser).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list client user err", true)
		return
	}
	resp.SuccessWithArrayData(c, clientUser, 0)
}

// GetClientUsers
// @Summary	get client user list
// @Tags	client-user
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	clientId	path	string	true	"client id"	default(default)
// @Param	subId		path	string	true	"subId"
// @Success	200
// @Router	/accounts/admin/{tenant}/clients/{clientId}/users/{subId} [get]
func GetClientUsers(c *gin.Context) {
	var clientUser struct {
		Sub      string `json:"sub"`
		ClientId string `json:"clientId"`
		model.User
	}
	clientId := c.Param("clientId")
	subId := c.Param("subId")
	if err := global.DB.Table("client_users cu").
		Select("cu.id, cu.sub sub, cu.client_id, u.username username, u.phone, u.email, u.first_name, u.last_name, u.display_name, u.role, u.avatar, u.from, u.meta").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ? AND cu.sub = ?", internal.GetTenant(c).Id, clientId, subId).
		Find(&clientUser).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "get client user err")
		return
	}

	if clientUser.Username == "" {
		resp.ErrorNotFound(c, "no such client user")
		return
	}
	resp.SuccessWithData(c, clientUser)
}

// UpdateUserMeta
// @Summary	update user
// @Tags	client-user
// @Param	tenant		path	string			true	"tenant"	default(default)
// @Param	clientId	path	string			true	"client id"	default(default)
// @Param	subId		path	string			true	"sub id"
// @Param	bd			body	string			true	"user body"
// @Success	200
// @Router	/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/meta [put]
func UpdateUserMeta(c *gin.Context) {
	meta, err := io.ReadAll(c.Request.Body)
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	user, err := service.GetUserBySubId(internal.GetTenant(c).Id, c.Param("clientId"), c.Param("subId"))
	if err != nil {
		resp.ErrorSqlSelect(c, err, "no such user")
		return
	}

	user.Meta = string(meta)
	if err = global.DB.Select("meta").Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update user meta err")
		return
	}
	resp.Success(c)
}

// UpdateUserPassword
// @Summary	update user
// @Tags	client-user
// @Param	tenant		path	string			true	"tenant"	default(default)
// @Param	clientId	path	string			true	"client id"	default(default)
// @Param	subId		path	string			true	"sub id"
// @Param	bd			body	ModifyPassword	true	"user body"
// @Success	200
// @Router	/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/password [put]
func UpdateUserPassword(c *gin.Context) {
	var u ModifyPassword
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	user, err := service.GetUserBySubId(internal.GetTenant(c).Id, c.Param("clientId"), c.Param("subId"))
	if err != nil {
		resp.ErrorSqlSelect(c, err, "no such user")
		return
	}

	// 检查旧密码
	if ok := utils.CheckPasswordHash(u.OldPassword, user.PasswordHash); !ok {
		resp.ErrorPassword(c, "password hash err")
		return
	}

	newHash, err := utils.HashPassword(u.NewPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}
	user.PasswordHash = newHash
	if err = global.DB.Select("password_hash").Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update user password err")
		return
	}
	resp.Success(c)
}

// UpdateUserProfile
// @Summary	update user
// @Tags	client-user
// @Param	tenant		path	string				true	"tenant"	default(default)
// @Param	clientId	path	string				true	"client id"	default(default)
// @Param	subId		path	string				true	"sub id"
// @Param	bd			body	dto.UserAdminDto	true	"user body"
// @Success	200
// @Router	/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/profile [put]
func UpdateUserProfile(c *gin.Context) {
	var u model.User
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	tenant := internal.GetTenant(c)
	user, err := service.GetUserBySubId(tenant.Id, c.Param("clientId"), c.Param("subId"))
	if err != nil {
		resp.ErrorSqlSelect(c, err, "service.GetUserBySubId err")
		return
	}
	user.Username = u.Username
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.DisplayName = u.DisplayName
	user.Email = u.Email
	user.Phone = u.Phone
	user.Avatar = u.Avatar
	if err = global.DB.Select("username", "first_name", "last_name", "display_name", "email", "phone", "avatar").
		Where("id = ?", user.Id).Updates(user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update tenant user err")
		return
	}

	resp.Success(c)
}

func AddClientUserRoute(rg *gin.RouterGroup) {
	rg.GET("/clients/:clientId/users", ListClientUsers)
	rg.GET("/clients/:clientId/users/:subId", GetClientUsers)
	rg.PUT("/clients/:clientId/users/:subId/meta", UpdateUserMeta)
	rg.PUT("/clients/:clientId/users/:subId/password", UpdateUserPassword)
	rg.PUT("/clients/:clientId/users/:subId/profile", UpdateUserProfile)
}

package authentication

import (
	"alfred/backend/endpoint/dto"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/middlewares"
	"alfred/backend/service"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) model.User {
	return c.MustGet("user").(model.User)
}

// GetUserDetail
// @Summary	get user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200	{object}	dto.UserDto
// @Router	/accounts/{tenant}/me [get]
func GetUserDetail(c *gin.Context) {
	user := GetUser(c)
	tenantName := "default"
	clientId := "default"

	tenantId, err := service.GetTenantIdByTenantName(tenantName)
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	sub, err := service.GetAlfredClientUser(clientId, tenantId, user.Id)
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	user.Sub = sub

	resp.SuccessWithData(c, user.Dto())
}

// UpdateUserDetail
// @Summary	update user
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Body	request	body	dto.UserProfileDto	true	"request"
// @Success	200	{object}	dto.UserDto
// @Router	/accounts/{tenant}/me [put]
func UpdateUserDetail(c *gin.Context) {
	user := GetUser(c)
	var u dto.UserDto
	if err := c.BindJSON(&u); err != nil {
		user.FirstName = u.FirstName
		user.LastName = u.DisplayName
		user.DisplayName = u.DisplayName
	}
	global.DB.Save(&user)
}

// GetUserProfile
// @Summary	get user profile
// @Tags	user
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200	{object}	dto.UserProfileDto
// @Router	/accounts/{tenant}/me/profile [get]
func GetUserProfile(c *gin.Context) {
	user := GetUser(c)
	resp.SuccessWithData(c, user.ProfileDto())
}

func AddUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", middlewares.Authorized(), GetUserDetail)
	rg.PUT("/me", middlewares.Authorized(), UpdateUserDetail)
	rg.GET("/me/profile", middlewares.Authorized(), GetUserProfile)
}

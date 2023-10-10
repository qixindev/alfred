package authentication

import (
	"alfred/internal/endpoint/dto"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/global"
	"alfred/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) *model.User {
	return c.MustGet("user").(*model.User)
}

// GetUserDetail godoc
//
//	@Summary	get user
//	@Schemes
//	@Description	get user
//	@Tags			user
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Success		200		{object}	dto.UserDto
//	@Router			/accounts/{tenant}/me [get]
func GetUserDetail(c *gin.Context) {
	user := GetUser(c)
	tenantName := "default"
	clientId := "default"

	tenantId, err := service.GetTenantIdByTenantName(tenantName)
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	sub, err := service.GetSub(clientId, tenantId, user.Id)
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	user.Sub = sub

	resp.SuccessWithData(c, user.Dto())

}

// UpdateUserDetail godoc
//
//	@Summary	update user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Body			request	body		dto.UserProfileDto	true	"request"
//	@Success		200		{object}	dto.UserDto
//	@Router			/accounts/{tenant}/me [put]
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

// GetUserProfile godoc
//
//	@Summary	get user profile
//	@Schemes
//	@Description	get user profile
//	@Tags			user
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Success		200		{object}	dto.UserProfileDto
//	@Router			/accounts/{tenant}/me/profile [get]
func GetUserProfile(c *gin.Context) {
	user := GetUser(c)
	resp.SuccessWithData(c, user.ProfileDto())
}

func AddUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", middlewares.Authorized(false), GetUserDetail)
	rg.PUT("/me", middlewares.Authorized(false), UpdateUserDetail)
	rg.GET("/me/profile", middlewares.Authorized(false), GetUserProfile)
}

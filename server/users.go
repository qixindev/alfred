package server

import (
	"accounts/global"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}

// GetUserDetail godoc
//
//	@Summary	get user
//	@Schemes
//	@Description	get user
//	@Tags			user
//	@Param			tenant	path		string	true	"tenant"
//	@Success		200		{object}	dto.UserDto
//	@Router			/accounts/{tenant}/me [get]
func GetUserDetail(c *gin.Context) {
	user := GetUser(c)
	c.JSON(http.StatusOK, user.Dto())
}

// UpdateUserDetail godoc
//
//	@Summary	update user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path		string	true	"tenant"
//	@Body			request	body							dto.UserProfileDto	true	"request"
//	@Success		200		{object}	dto.UserDto
//	@Router			/accounts/{tenant}/me [put]
func UpdateUserDetail(c *gin.Context) {
	user := GetUser(c)
	var u dto.UserDto
	err := c.BindJSON(&u)
	if err != nil {
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
//	@Param			tenant	path		string	true	"tenant"
//	@Success		200		{object}	dto.UserProfileDto
//	@Router			/accounts/{tenant}/me/profile [get]
func GetUserProfile(c *gin.Context) {
	user := GetUser(c)
	c.JSON(http.StatusOK, user.ProfileDto())
}

func addUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", middlewares.Authorized(false), GetUserDetail)
	rg.PUT("/me", middlewares.Authorized(false), UpdateUserDetail)
	rg.GET("/me/profile", middlewares.Authorized(false), GetUserProfile)
}

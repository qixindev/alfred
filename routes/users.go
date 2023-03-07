package routes

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/models/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}

func addUsersRoutes(rg *gin.RouterGroup) {
	rg.Use(middlewares.Authorized)
	rg.GET("/me", func(c *gin.Context) {
		user := GetUser(c)
		c.JSON(http.StatusOK, user.Dto())
	})

	rg.PUT("/me", func(c *gin.Context) {
		user := GetUser(c)
		var u dto.UserDto
		err := c.BindJSON(&u)
		if err != nil {
			user.FirstName = u.FirstName
			user.LastName = u.DisplayName
			user.DisplayName = u.DisplayName
		}
		data.DB.Save(&user)
	})

	rg.GET("/me/profile", func(c *gin.Context) {
		user := GetUser(c)
		c.JSON(http.StatusOK, user.ProfileDto())
	})

	rg.PUT("/me/profile", func(c *gin.Context) {
		user := GetUser(c)
		var u dto.UserProfileDto
		err := c.BindJSON(&u)
		if err != nil {
			user.FirstName = u.FirstName
			user.LastName = u.DisplayName
			user.DisplayName = u.DisplayName
		}
		data.DB.Save(&user)
	})
}

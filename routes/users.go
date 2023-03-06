package routes

import (
	"accounts/data"
	"accounts/models"
	"accounts/models/dto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}

func Authorized(c *gin.Context) {
	tenant := GetTenant(c)
	session := sessions.Default(c)
	tenantName := session.Get("tenant")
	if tenant.Name != tenantName {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	username := session.Get("user")
	var user models.User
	if data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, username).Error != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user", &user)
	c.Next()
}

func addUsersRoutes(rg *gin.RouterGroup) {
	rg.Use(Authorized)
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

package routes

import (
	"accounts/data"
	"accounts/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func addLoginRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", func(c *gin.Context) {
		login := c.PostForm("login")
		password := c.PostForm("password")

		if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		tenant := GetTenant(c)

		var user models.User
		if data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, login).Error != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		if checkPasswordHash(password, user.PasswordHash) == false {
			c.Status(http.StatusUnauthorized)
			return
		}

		session := sessions.Default(c)
		session.Set("tenant", tenant.Name)
		session.Set("user", user.Username)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
	})

	rg.GET("/logout", Authorized, func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Status(http.StatusBadRequest)
			return
		}
		session.Delete("tenant")
		session.Delete("user")
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
	})

	rg.POST("/register", func(c *gin.Context) {
		tenant := GetTenant(c)
		login := c.PostForm("login")
		password := c.PostForm("password")

		var user models.User
		err := data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, login).Error
		if err == nil {
			c.Status(http.StatusConflict)
			return
		}

		hash, err := hashPassword(password)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		newUser := models.User{
			TenantId:     tenant.Id,
			Username:     login,
			PasswordHash: hash,
		}
		if err := data.DB.Create(&newUser).Error; err != nil {
			log.Print(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	})
}

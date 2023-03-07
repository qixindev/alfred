package routes

import (
	"accounts/auth"
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		tenant := middlewares.GetTenant(c)

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

	rg.GET("/logout", middlewares.Authorized, func(c *gin.Context) {
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
		tenant := middlewares.GetTenant(c)
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

	rg.GET("/login/:provider", func(c *gin.Context) {
		providerName := c.Param("provider")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "name = ?", providerName).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var userInfo *auth.UserInfo = nil
		if provider.Type == "oauth2" {
			var err error
			userInfo, err = auth.GetOAuth2User(provider)
			if err != nil {
				c.Status(http.StatusUnauthorized)
			}
		}

		var providerUser models.ProviderUser
		var existingUser *models.User = nil
		if middlewares.TenantDB(c).First(&providerUser, "provider_id = ? AND name = ?", provider.Id, userInfo.Name).Error != nil {
			// If logged in, bind to current user.
			user, err := middlewares.GetUserStandalone(c)
			if err != nil {
				user := &models.User{
					Username:         uuid.NewString(),
					FirstName:        userInfo.FirstName,
					LastName:         userInfo.LastName,
					DisplayName:      userInfo.DisplayName,
					Email:            userInfo.Email,
					EmailVerified:    false,
					Phone:            userInfo.Phone,
					PhoneVerified:    false,
					TwoFactorEnabled: false,
					Disabled:         false,
					TenantId:         provider.TenantId,
				}
				if data.DB.Create(user).Error != nil {
					c.Status(http.StatusConflict)
					return
				}
			}

			providerUser.TenantId = provider.TenantId
			providerUser.ProviderId = provider.Id
			providerUser.UserId = user.Id
			provider.Name = userInfo.Name
			if data.DB.Create(&provider).Error != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			existingUser = user
		} else {
			if middlewares.TenantDB(c).First(existingUser, "id = ?", providerUser.UserId).Error != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		session := sessions.Default(c)
		tenant := middlewares.GetTenant(c)
		session.Set("tenant", tenant.Name)
		session.Set("user", existingUser.Username)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
	})
}

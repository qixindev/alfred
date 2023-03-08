package routes

import (
	"accounts/auth"
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/utils"
	"errors"
	"fmt"
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

func getAuthProvider(tenantId uint, providerName string) (auth.AuthProvider, error) {
	var provider models.Provider
	if err := data.DB.First(&provider, "tenant_id = ? AND name = ?", tenantId, providerName).Error; err != nil {
		return nil, err
	}
	if provider.Type == "oauth2" {
		var providerOAuth2 auth.ProviderOAuth2
		if err := data.DB.First(&providerOAuth2, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, err
		}
		return providerOAuth2, nil
	}
	if provider.Type == "dingtalk" {
		var providerDingTalk auth.ProviderDingTalk
		if err := data.DB.First(&providerDingTalk, "tenant_id = ? AND provider_id = ?", tenantId, provider.Id).Error; err != nil {
			return nil, err
		}
		return providerDingTalk, nil
	}
	return nil, errors.New("provider config not found")
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

	rg.GET("/login/:provider", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		providerName := c.Param("provider")
		authProvider, err := getAuthProvider(tenant.Id, providerName)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenant.Name, providerName)
		location := authProvider.Auth(redirectUri)
		c.Redirect(http.StatusFound, location)
	})

	rg.GET("/providers", func(c *gin.Context) {
		var providers []models.Provider
		if middlewares.TenantDB(c).Find(&providers).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(providers, models.Provider2Dto))
	})
	rg.GET("/providers/:provider", func(c *gin.Context) {
		providerName := c.Param("provider")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "name = ?", providerName).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, provider.Dto())
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

	rg.GET("/logged-in/:provider", func(c *gin.Context) {
		providerName := c.Param("provider")
		var provider models.Provider
		if middlewares.TenantDB(c).First(&provider, "name = ?", providerName).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}

		authProvider, err := getAuthProvider(provider.TenantId, provider.Name)
		if err != nil {
			c.String(http.StatusNotFound, "provider not found")
		}
		userInfo, err := authProvider.Login(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		var providerUser models.ProviderUser
		existingUser := &models.User{}
		if middlewares.TenantDB(c).First(&providerUser, "provider_id = ? AND name = ?", provider.Id, userInfo.Sub).Error != nil {
			// Current bind not found.
			// If logged in, bind to current user.
			user, err := middlewares.GetUserStandalone(c)
			if err != nil {
				// If not logged in, create new user.
				newUser := models.User{
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
				if data.DB.Create(&newUser).Error != nil {
					c.Status(http.StatusConflict)
					return
				}
				user = &newUser
			}

			providerUser.TenantId = provider.TenantId
			providerUser.ProviderId = provider.Id
			providerUser.UserId = user.Id
			providerUser.Name = userInfo.Sub
			if err := data.DB.Create(&providerUser).Error; err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			existingUser = user
		} else {
			if err := middlewares.TenantDB(c).First(existingUser, "id = ?", providerUser.UserId).Error; err != nil {
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

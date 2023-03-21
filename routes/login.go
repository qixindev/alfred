package routes

import (
	"accounts/auth"
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/utils"
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

// Login godoc
//
//	@Summary	login a user
//	@Schemes
//	@Description	login using username and password
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"
//	@Param			login		formData	string	true	"username"
//	@Param			password	formData	string	true	"password"
//	@Param			next		query		string	false	"next"
//	@Success		302
//	@Router			/accounts/{tenant}/login [post]
func Login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	tenant := middlewares.GetTenant(c)

	var user models.User
	if err := data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, login).Error; err != nil {
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
	session.Set("userId", user.Id)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	next := c.Query("next")
	if next == "" {
		next = c.PostForm("next")
	}
	if next != "" {
		c.Redirect(http.StatusFound, next)
		return
	}
}

// LoginToProvider godoc
//
//	@Summary	login via a provider
//	@Schemes
//	@Description	login via a provider
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"
//	@Param			provider	path	string	true	"provider"
//	@Param			next		query	string	false	"next"
//	@Success		302
//	@Router			/accounts/{tenant}/login/{provider} [get]
func LoginToProvider(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	providerName := c.Param("provider")
	authProvider, err := auth.GetAuthProvider(tenant.Id, providerName)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	redirectUri := fmt.Sprintf("%s/%s/logged-in/%s", utils.GetHostWithScheme(c), tenant.Name, providerName)
	location, err := authProvider.Auth(redirectUri)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	next := c.Query("next")
	if next != "" {
		session := sessions.Default(c)
		session.Set("next", next)
		session.Save()
	}
	c.Redirect(http.StatusFound, location)
}

// ListProviders godoc
//
//	@Summary	List all providers
//	@Schemes
//	@Description	list login providers
//	@Tags			login
//	@Param			tenant	path		string	true	"tenant"
//	@Success		200		{object}	[]dto.ProviderDto
//	@Router			/accounts/{tenant}/login/providers [get]
func ListProviders(c *gin.Context) {
	var providers []models.Provider
	if middlewares.TenantDB(c).Find(&providers).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(providers, models.Provider2Dto))
}

// GetProvider godoc
//
//	@Summary	get a provider
//	@Schemes
//	@Description	get a login provider
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"
//	@Param			provider	path		string	true	"provider"
//	@Success		200			{object}	dto.ProviderDto
//	@Router			/accounts/{tenant}/login/providers/{provider} [get]
func GetProvider(c *gin.Context) {
	providerName := c.Param("provider")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "name = ?", providerName).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, provider.Dto())
}

// Register godoc
//
//	@Summary	register a user
//	@Schemes
//	@Description	register using username and password
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"
//	@Param			login		formData	string	true	"username"
//	@Param			password	formData	string	true	"password"
//	@Success		200
//	@Router			/accounts/{tenant}/register [post]
func Register(c *gin.Context) {
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
}

// Logout godoc
//
//	@Summary	logout current user
//	@Schemes
//	@Description	logout current user
//	@Tags			login
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/{tenant}/logout [get]
func Logout(c *gin.Context) {
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
}

// ProviderCallback godoc
//
//	@Summary	provider callback
//	@Schemes
//	@Description	provider callback
//	@Tags			login
//	@Param			tenant		path	string	true	"tenant"
//	@Param			provider	path	string	true	"provider"
//	@Param			code		query	string	true	"code"
//	@Success		302
//	@Success		200
//	@Router			/accounts/{tenant}/logged-in/{provider} [get]
func ProviderCallback(c *gin.Context) {
	providerName := c.Param("provider")
	var provider models.Provider
	if middlewares.TenantDB(c).First(&provider, "name = ?", providerName).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	authProvider, err := auth.GetAuthProvider(provider.TenantId, provider.Name)
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
	next := utils.GetString(session.Get("next"))
	session.Delete("next")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	if next != "" {
		c.Redirect(http.StatusFound, next)
		return
	}
}

func addLoginRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", Login)
	rg.GET("/login/:provider", LoginToProvider)
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:provider", GetProvider)
	rg.GET("/logout", middlewares.Authorized(false), Logout)
	rg.POST("/register", Register)
	rg.GET("/logged-in/:provider", ProviderCallback)
}

package authentication

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/global"
	"accounts/pkg/middlewares"
	"accounts/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Login godoc
//
//	@Summary	login a user
//	@Schemes
//	@Description	login using username and password
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			login		formData	string	true	"username"
//	@Param			password	formData	string	true	"password"
//	@Param			next		query		string	false	"next"
//	@Success		302
//	@Router			/accounts/{tenant}/login [post]
func Login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
		resp.ErrorRequestWithMsg(c, nil, "username or password should not be empty")
		return
	}

	tenant := internal.GetTenant(c)

	var user model.User
	if err := global.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, login).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}

	if utils.CheckPasswordHash(password, user.PasswordHash) == false {
		resp.ErrorUnauthorized(c, nil, "incorrect password")
		return
	}

	session := sessions.Default(c)
	session.Set("tenant", tenant.Name)
	session.Set("user", user.Username)
	session.Set("userId", user.Id)
	if err := session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
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

// Register godoc
//
//	@Summary	register a user
//	@Schemes
//	@Description	register using username and password
//	@Tags			login
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			login		formData	string	true	"username"
//	@Param			password	formData	string	true	"password"
//	@Success		200
//	@Router			/accounts/{tenant}/register [post]
func Register(c *gin.Context) {
	tenant := internal.GetTenant(c)
	login := c.PostForm("login")
	password := c.PostForm("password")
	var user model.User
	if err := global.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, login).Error; err == nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		resp.ErrorUnauthorized(c, nil, "hashPassword err")
		return
	}

	newUser := model.User{
		TenantId:     tenant.Id,
		Username:     login,
		PasswordHash: hash,
	}
	if err = global.DB.Create(&newUser).Error; err != nil {
		resp.ErrorUnknown(c, err, "create user err")
		return
	}
	resp.Success(c)
}

// Logout godoc
//
//	@Summary	logout current user
//	@Schemes
//	@Description	logout current user
//	@Tags			login
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/{tenant}/logout [get]
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		resp.ErrorUnknown(c, nil, "session user is nil")
		return
	}
	session.Delete("tenant")
	session.Delete("user")
	if err := session.Save(); err != nil {
		resp.ErrorSaveSession(c, err)
		return
	}
	resp.Success(c)
}

func AddLoginRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", Login)
	rg.GET("/login/:provider", LoginToProvider)
	rg.GET("/providers", ListProviders)
	rg.GET("/providers/:provider", GetProvider)
	rg.GET("/logout", middlewares.Authorized(false), Logout)
	rg.POST("/register", Register)
	rg.GET("/logged-in/:provider", ProviderCallback)
}

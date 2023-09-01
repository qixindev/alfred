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
		resp.ErrorRequestWithMsg(c, "username or password should not be empty")
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

// GetLoginProtocol godoc
//
//	@Summary	logout current user
//	@Schemes
//	@Description	logout current user
//	@Tags			login
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/{tenant}/login/proto [get]
func GetLoginProtocol(c *gin.Context) {
	resp.SuccessWithArrayData(c, []gin.H{
		{
			"name": "用户服务协议",
			"url":  "https://devbackup.blob.core.chinacloudapi.cn/picture/%E7%94%A8%E6%88%B7%E6%9C%8D%E5%8A%A1%E5%8D%8F%E8%AE%AE.md?sp=r&st=2023-08-31T07:30:43Z&se=2123-08-31T15:30:43Z&spr=https&sv=2022-11-02&sr=b&sig=PsaLy37JxY6Buhp9z9fwv6PiofeKbCwGzZNv7iI2bg4%3D",
		}, {
			"name": "隐私保护声明",
			"url":  "https://devbackup.blob.core.chinacloudapi.cn/picture/%E9%9A%90%E7%A7%81%E4%BF%9D%E6%8A%A4%E5%A3%B0%E6%98%8E.md?sp=r&st=2023-08-31T07:24:53Z&se=2123-08-31T15:24:53Z&spr=https&sv=2022-11-02&sr=b&sig=Jpd16wFlwtpF7GNry%2FsH8vDQVpoDWOERRjm273QGR7I%3D",
		},
	}, 2)
}

func AddLoginRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", Login)                    // 账号密码登录
	rg.GET("/login/proto", GetLoginProtocol)    // 获取登录隐私协议
	rg.GET("/login/:provider", LoginToProvider) // 第三方登录重定向
	rg.GET("/providers", ListProviders)         // 第三方信息
	rg.GET("/providers/:provider", GetProvider) // 第三方具体信息
	rg.GET("/logout", middlewares.Authorized(false), Logout)
	rg.POST("/register", Register)                   // 注册
	rg.GET("/logged-in/:provider", ProviderCallback) // 验证第三方登录是否成功
}

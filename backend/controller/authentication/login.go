package authentication

import (
	"alfred/backend/controller/internal"
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/middlewares"
	"alfred/backend/pkg/utils"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

// Login
// @Summary	login using username and password
// @Tags	login
// @Param	tenant		path		string	true	"tenant"	default(default)
// @Param	login		formData	string	true	"username"
// @Param	password	formData	string	true	"password"
// @Param	next		query		string	false	"next"
// @Success	302
// @Router	/accounts/{tenant}/login [post]
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

// Register
// @Summary	register using username and password
// @Tags	login
// @Param	tenant		path		string	true	"tenant"	default(default)
// @Param	login		formData	string	true	"username"
// @Param	password	formData	string	true	"password"
// @Success	200
// @Router	/accounts/{tenant}/register [post]
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
		From:         "register",
		Meta:         "{}",
	}
	if err = global.DB.Create(&newUser).Error; err != nil {
		resp.ErrorUnknown(c, err, "create user err")
		return
	}
	resp.Success(c)
}

// Logout
// @Summary	logout current user
// @Tags	login
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200
// @Router	/accounts/{tenant}/logout [get]
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

// GetLoginProtocol
// @Summary	logout current user
// @Tags	login
// @Param	tenant		path	string	true	"tenant"	default(default)
// @Param	fileName	path	string	true	"fileName"	default(default)
// @Success	200
// @Router	/accounts/{tenant}/login/proto/{fileName} [get]
func GetLoginProtocol(c *gin.Context) {
	fileName := c.Param("fileName")

	filePaths := map[string]string{
		"userServiceAgreement": "docs/用户服务协议.md",
		"privacyStatement":     "docs/隐私保护声明.md",
	}

	filePath, exists := filePaths[fileName]
	if !exists {
		resp.ErrorUnknown(c, errors.New("invalid fileName"), "无效的文件名")
		return
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		resp.ErrorUnknown(c, err, "无法读取文件")
		return
	}

	resp.SuccessWithData(c, gin.H{
		"name":    fileName,
		"content": string(fileContent),
	})
}

func AddLoginRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", Login)                              // 账号密码登录
	rg.GET("/login/proto/:fileName", GetLoginProtocol)    // 获取登录协议
	rg.GET("/providers", ListProviders)                   // 第三方信息列表
	rg.GET("/providers/:provider", GetProvider)           // 第三方具体信息
	rg.GET("/providers/:provider/login", LoginToProvider) // 第三方登录重定向
	rg.GET("/logout", middlewares.Authorized(false), Logout)
	rg.POST("/register", Register) // 注册
}

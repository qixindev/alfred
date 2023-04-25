package middlewares

import (
	"accounts/global"
	"accounts/models"
	"accounts/utils"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"net/url"
	"strings"
)

func getTenant(c *gin.Context) *models.Tenant {
	return c.MustGet("tenant").(*models.Tenant)
}

func MultiTenancy(c *gin.Context) {
	tenantName := c.Param("tenant")
	fmt.Println(c.Request.URL.String())
	if tenantName == "" && !strings.HasPrefix(c.Request.URL.String(), "/accounts/admin/tenants") {
		c.AbortWithStatusJSON(http.StatusBadRequest, &gin.H{"message": "tenant should not be null"})
	}

	var tenant models.Tenant
	if global.DB.First(&tenant, "name = ?", tenantName).Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}
	tenantName = c.Request.Host
	if global.DB.First(&tenant, "name = ?", tenantName).Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}

	if global.DB.First(&tenant, "name = ?", "default").Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}

	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Tenant not found."})
	return
}

func GetUserStandalone(c *gin.Context) (*models.User, error) {
	tenant := getTenant(c)
	session := sessions.Default(c)
	tenantName := session.Get("tenant")
	if tenant.Name != tenantName {
		global.LOG.Error(fmt.Sprintf("tenant name err: %s %s", tenant.Name, tenantName))
		return nil, errors.New("")
	}

	var user models.User
	username := session.Get("user")
	if err := global.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, username).Error; err != nil {
		global.LOG.Error("get tenant user err: " + err.Error())
		return nil, errors.New("")
	}
	return &user, nil
}

func AuthorizedAdmin(c *gin.Context) {
	if c.GetHeader("Authorization") != "" {
		AuthAccessToken(c)
		return
	}

	user, err := GetUserStandalone(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "用户与租户不匹配")
		return
	}
	if user.Role != "owner" && user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "非管理员无权访问")
		return
	}
	c.Set("user", user)
	c.Next()
}

func AuthAccessToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	keys, err := utils.LoadRsaPrivateKeys("default")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "load key err"})
		global.LOG.Error("get private key err")
		return
	}

	var key *rsa.PrivateKey
	for _, key = range keys {
		claim := jwt.New(jwt.SigningMethodRS256)
		token, err := jwt.ParseWithClaims(tokenString, claim.Claims, func(token *jwt.Token) (interface{}, error) {
			return key.Public(), nil
		})

		if err == nil && token.Valid {
			return
		}
		global.LOG.Warn(fmt.Sprintf("%s token valid err: %s", "default", err))
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token invalidate"})
}

func Authorized(redirectToLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserStandalone(c)
		if err != nil {
			if redirectToLogin {
				t := getTenant(c)
				h := utils.GetHostWithScheme(c)
				base := fmt.Sprintf("%s/%s", h, t.Name)
				next := fmt.Sprintf("%s/oauth2/auth", base)
				location := fmt.Sprintf("%s/login?next=%s", base, url.QueryEscape(next))
				c.Redirect(http.StatusFound, location)
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				c.Abort()
			}
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

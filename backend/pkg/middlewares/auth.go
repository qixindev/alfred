package middlewares

import (
	"alfred/backend/endpoint/resp"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"crypto/rsa"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func getTenant(c *gin.Context) *model.Tenant {
	return c.MustGet("tenant").(*model.Tenant)
}

func MultiTenancy(c *gin.Context) {
	tenantName := c.Param("tenant")
	if strings.HasPrefix(c.Request.URL.String(), "/accounts/admin/tenants") {
		tenantName = "default"
	}
	if tenantName == "" {
		resp.ErrorRequest(c, nil)
		return
	}

	var tenant model.Tenant
	if global.DB.First(&tenant, "name = ?", tenantName).Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}

	resp.ErrorNotFound(c, "Tenant not found.")
}

func AuthorizedAdmin(c *gin.Context) {
	if c.GetHeader("Authorization") != "" {
		AuthAccessToken(c)
		return
	}

	username := sessions.Default(c).Get("user")
	tenantName := c.Param("tenant")
	if username == nil {
		resp.ErrorNotLogin(c, nil)
		return
	}
	if tenantName == "" {
		return
	}

	var user model.User
	if err := global.DB.Table("users as u").Select("u.username, u.role, u.phone, u.email").
		Joins("LEFT JOIN tenants as t ON t.id = u.tenant_id").
		First(&user, "t.name = ? AND u.username = ?", tenantName, username).Error; err != nil {
		resp.ErrorForbidden(c, err, "get tenant user err")
		return
	}

	if user.Role != "owner" && user.Role != "admin" {
		resp.ErrorUnauthorized(c, nil, "非管理员无权访问")
		return
	}
	c.Set("user", user)
}

func AuthAccessToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	keys, err := utils.LoadRsaPrivateKeys("default")
	if err != nil {
		resp.ErrorUnauthorized(c, err, "load private key err")
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

	resp.ErrorUnauthorized(c, nil, "token invalidate")
}

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		tenant := getTenant(c)
		session := sessions.Default(c)
		username := session.Get("user")
		if err := global.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, username).Error; err != nil {
			resp.ErrorNotLogin(c, err)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

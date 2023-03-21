package middlewares

import (
	"accounts/data"
	"accounts/models"
	"accounts/utils"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
)

func TenantDB(c *gin.Context) *gorm.DB {
	tenant := GetTenant(c)
	return data.WithTenant(tenant.Id)
}

func GetTenant(c *gin.Context) *models.Tenant {
	return c.MustGet("tenant").(*models.Tenant)
}

func MultiTenancy(c *gin.Context) {
	tenantName := c.Param("tenant")
	var tenant models.Tenant
	if data.DB.First(&tenant, "name = ?", tenantName).Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}
	tenantName = c.Request.Host
	if data.DB.First(&tenant, "name = ?", tenantName).Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}

	if data.DB.First(&tenant, "name = ?", "default").Error == nil {
		c.Set("tenant", &tenant)
		c.Next()
		return
	}

	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Tenant not found."})
	return
}

func GetUserStandalone(c *gin.Context) (*models.User, error) {
	tenant := GetTenant(c)
	session := sessions.Default(c)
	tenantName := session.Get("tenant")
	if tenant.Name != tenantName {
		fmt.Printf("tenant name err: %s %s\n", tenant.Name, tenantName)
		return nil, errors.New("")
	}
	username := session.Get("user")
	var user models.User
	if data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, username).Error != nil {
		return nil, errors.New("")
	}
	return &user, nil
}

func AuthorizedAdmin(c *gin.Context) {
	user, err := GetUserStandalone(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "租户不匹配")
		return
	}
	if user.Role != "owner" && user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "非管理员无权访问")
		return
	}
	c.Set("user", user)
	c.Next()
}

func Authorized(redirectToLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserStandalone(c)
		if err != nil {
			if redirectToLogin {
				t := GetTenant(c)
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

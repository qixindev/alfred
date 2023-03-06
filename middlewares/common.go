package middlewares

import (
	"accounts/data"
	"accounts/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func TenantDB(c *gin.Context) *gorm.DB {
	tenant := GetTenant(c)
	return data.DB.Where("tenant_id = ?", tenant.Id)
}

func GetTenant(c *gin.Context) *models.Tenant {
	return c.MustGet("tenant").(*models.Tenant)
}

func MultiTenancy(c *gin.Context) {
	tenantName := c.Param("tenant")
	var tenant models.Tenant
	if data.DB.First(&tenant, "name = ?", tenantName).Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Tenant not found."})
		return
	}
	c.Set("tenant", &tenant)
	c.Next()
}

func AuthorizedAdmin(c *gin.Context) {
	tenant := GetTenant(c)
	session := sessions.Default(c)
	tenantName := session.Get("tenant")
	if tenant.Name != tenantName {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	username := session.Get("user")
	var user models.User
	if data.DB.First(&user, "tenant_id = ? AND username = ?", tenant.Id, username).Error != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Role != "owner" && user.Role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user", &user)
	c.Next()
}

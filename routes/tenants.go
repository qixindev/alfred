package routes

import (
	"accounts/data"
	"accounts/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

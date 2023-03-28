package internal

import (
	"accounts/global"
	"accounts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantDB(c *gin.Context) *gorm.DB {
	tenant := GetTenant(c)
	return global.DB.Where("tenant_id = ?", tenant.Id)
}

func GetTenant(c *gin.Context) *models.Tenant {
	return c.MustGet("tenant").(*models.Tenant)
}

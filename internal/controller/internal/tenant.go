package internal

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantDB(c *gin.Context) *gorm.DB {
	tenant := GetTenant(c)
	return global.DB.Where("tenant_id = ?", tenant.Id)
}

func GetTenant(c *gin.Context) *model.Tenant {
	return c.MustGet("tenant").(*model.Tenant)
}

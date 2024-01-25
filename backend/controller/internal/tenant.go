package internal

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
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

func (a *Api) SetTenant(tenant *model.Tenant) *Api {
	if a.c == nil {
		return a.setError(errors.New("gin context should not be nil"))
	}
	t, ok := a.c.Get("tenant")
	if !ok {
		return a.setError(errors.New("failed to get tenant from context"))
	}
	tenant, ok = t.(*model.Tenant)
	if !ok {
		return a
	}
	return a
}

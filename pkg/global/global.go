package global

import (
	"alfred/pkg/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	CONFIG *config.Config
	LOG    *zap.Logger
	DB     *gorm.DB
)

func WithTenant(tenantId uint) *gorm.DB {
	return DB.Where("tenant_id = ?", tenantId)
}

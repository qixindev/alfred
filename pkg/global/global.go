package global

import (
	"alfred/pkg/config"
	"github.com/allegro/bigcache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	CONFIG    *config.Config
	LOG       *zap.Logger
	DB        *gorm.DB
	CodeCache *bigcache.BigCache
)

func WithTenant(tenantId uint) *gorm.DB {
	return DB.Where("tenant_id = ?", tenantId)
}

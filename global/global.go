package global

import (
	"go.uber.org/zap"
)

var (
	CONFIG *config.SafeConfig
	LOG    *zap.Logger
)

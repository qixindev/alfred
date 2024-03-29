package service

import (
	"alfred/backend/pkg/global"
	"errors"
)

func deleteSource(tenantId uint, relayList []any, id any, name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if tenantId == 0 {
		return errors.New("tenantId cannot be negative")
	}

	for _, v := range relayList {
		if err := global.DB.Model(v).Where(name+" = ? AND tenant_id = ?", id, tenantId).Delete(&v).Error; err != nil {
			return err
		}
	}
	return nil
}

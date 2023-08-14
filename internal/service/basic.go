package service

import (
	"accounts/pkg/global"
)

func deleteSource(md any, relayList []any, id any, name string) error {
	for _, v := range relayList {
		if err := global.DB.Model(v).Where(name+" = ?", id).Error; err != nil {
			return err
		}
	}
	if err := global.DB.Model(md).Where("id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

package service

import (
	"accounts/internal/endpoint/req"
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"
)

func GetSmsModel(t string) (any, error) {
	switch t {
	case "tcloud":
		return model.SmsTcloud{}, nil
	}
	return nil, errors.New("no such type")
}

func GetSmsConfig(tenantId uint, connId uint, t string) (any, error) {
	md, err := GetSmsModel(t)
	if err != nil {
		return nil, err
	}
	if err = global.DB.Model(md).Where("tenant_id = ? AND sms_connector_id = ?", tenantId, connId).
		First(&md).Error; err != nil {
		return nil, err
	}

	return &md, nil
}

func CreateSmsConfig(t string, sms req.Sms) error {
	_, err := GetSmsModel(t)
	if err != nil {
		return err
	}
	if err = global.DB.Create(model.SmsConnector{
		TenantId: sms.TenantId,
		Name:     sms.Name,
		Type:     t,
	}).Error; err != nil {
		return err
	}

	// todo: save to sms

	return nil
}

func UpdateSmsConfig(tenantId uint, connId uint, t string, sms req.Sms) error {
	_, err := GetSmsModel(t)
	if err != nil {
		return err
	}
	if err = global.DB.Where("tenant_id = ? AND sms_connector_id = ?", tenantId, connId).Updates(model.SmsConnector{
		TenantId: sms.TenantId,
		Name:     sms.Name,
		Type:     t,
	}).Error; err != nil {
		return err
	}

	// todo: save to sms

	return nil
}

func DeleteSmsConfig(tenantId uint, connId uint, t string) error {
	md, err := GetSmsModel(t)
	if err != nil {
		return err
	}
	if err = global.DB.Model(md).Where("tenant_id = ? AND sms_connector_id = ?", tenantId, connId).
		Delete(&md).Error; err != nil {
		return err
	}

	if err = global.DB.Where("tenant_id = ? AND id = ?", tenantId, connId).
		Delete(model.SmsConnector{}).Error; err != nil {
		return err
	}

	return nil
}

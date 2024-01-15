package service

import (
	"alfred/backend/endpoint/req"
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"errors"
)

func GetSmsModel(t string) (model.InterfaceSms, error) {
	switch t {
	case "tcloud":
		return &model.SmsTcloud{}, nil
	case "alibaba":
		return &model.SmsAlibaba{}, nil
	}
	return nil, errors.New("no such type")
}

func GetSmsConfig(tenantId uint, connId uint, t string) (any, error) {
	md, err := GetSmsModel(t)
	if err != nil {
		return nil, err
	}
	if err = global.DB.Model(md).Where("tenant_id = ? AND sms_connector_id = ?", tenantId, connId).
		Preload("SmsConnector").First(md).Error; err != nil {
		return nil, err
	}

	return &md, nil
}

func CreateSmsConfig(t string, sms req.Sms) (*model.SmsConnector, error) {
	md, err := GetSmsModel(t)
	if err != nil {
		return nil, err
	}

	conn := model.SmsConnector{
		TenantId: sms.TenantId,
		Name:     sms.Name,
		Type:     t,
	}
	if err = global.DB.Create(&conn).Error; err != nil {
		return nil, err
	}

	sms.Id = conn.Id
	if err = global.DB.Create(md.Save(sms)).Error; err != nil {
		return nil, err
	}
	return &conn, nil
}

func UpdateSmsConfig(tenantId uint, connId uint, t string, sms req.Sms) error {
	md, err := GetSmsModel(t)
	if err != nil {
		return err
	}
	if err = global.DB.Where("tenant_id = ? AND id = ?", tenantId, connId).Updates(model.SmsConnector{
		TenantId: sms.TenantId,
		Name:     sms.Name,
		Type:     t,
	}).Error; err != nil {
		return err
	}

	if err = global.DB.Where("tenant_id = ? AND sms_connector_id = ?", tenantId, connId).
		Updates(md.Save(sms)).Error; err != nil {
		return err
	}

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

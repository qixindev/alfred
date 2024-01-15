package sms

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
)

type Connector interface {
	Send(number string, contents string) error
}

func GetConnector(tenantId, connectorId uint) (Connector, error) {
	var c model.SmsConnector
	if err := global.DB.First(&c, "tenant_id = ? AND id = ?", tenantId, connectorId).Error; err != nil {
		return nil, err
	}
	var connector Connector
	if c.Type == "tcloud" {
		var config model.SmsTcloud
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		smsConfig := Tcloud{Config: config}
		connector = &smsConfig
	} else if c.Type == "alibaba" {
		var config model.SmsAlibaba
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		smsConfig := Alibaba{Config: config}
		connector = &smsConfig
	}
	return connector, nil
}

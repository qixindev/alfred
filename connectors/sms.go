package connectors

import (
	"accounts/global"
	"accounts/models"
)

type SmsConnector interface {
	Send(number string, contents []string) error
}

func GetConnector(tenantId, connectorId uint) (SmsConnector, error) {
	var c models.SmsConnector
	if err := global.DB.First(&c, "tenant_id = ? AND id = ?", tenantId, connectorId).Error; err != nil {
		return nil, err
	}
	var connector SmsConnector
	if c.Type == "tcloud" {
		var config models.SmsTcloud
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		sms := SmsTcloud{Config: config}
		connector = &sms
	}
	return connector, nil
}

package connectors

import (
	"accounts/data"
	"accounts/models"
)

type SmsConnector interface {
	Send(number, contents []string) error
}

func GetService(c models.SmsConnector) (SmsConnector, error) {
	var service SmsConnector
	if c.Type == "tcloud" {
		var config models.SmsTcloud
		if err := data.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		sms := SmsTcloud{Config: config}
		service = &sms
	}
	return service, nil
}

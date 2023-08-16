package service

import (
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SmsTcloud struct {
	Config model.SmsTcloud
}

func (s *SmsTcloud) Send(number string, contents []string) error {
	credential := common.NewCredential(s.Config.SecretId, s.Config.SecretKey)
	client, _ := sms.NewClient(credential, s.Config.Region, profile.NewClientProfile())
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.Config.SdkAppId)
	request.SignName = common.StringPtr(s.Config.SignName)
	request.TemplateId = common.StringPtr(s.Config.TemplateId)
	request.TemplateParamSet = common.StringPtrs(contents)
	request.PhoneNumberSet = common.StringPtrs([]string{number})

	if _, err := client.SendSms(request); err != nil {
		return err
	}
	return nil
}

type SmsConnector interface {
	Send(number string, contents []string) error
}

func GetConnector(tenantId, connectorId uint) (SmsConnector, error) {
	var c model.SmsConnector
	if err := global.DB.First(&c, "tenant_id = ? AND id = ?", tenantId, connectorId).Error; err != nil {
		return nil, err
	}
	var connector SmsConnector
	if c.Type == "tcloud" {
		var config model.SmsTcloud
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		smsConfig := SmsTcloud{Config: config}
		connector = &smsConfig
	}
	return connector, nil
}

func GetConnectorDetails(c model.SmsConnector) (any, error) {
	switch c.Type {
	case "tcloud":
		var config model.SmsTcloud
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		return config, nil
	}
	return nil, errors.New("no such connector")
}

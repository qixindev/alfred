package connectors

import (
	"accounts/models"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SmsTcloud struct {
	Config models.SmsTcloud
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

	_, err := client.SendSms(request)
	if err != nil {
		return err
	}
	return nil
}

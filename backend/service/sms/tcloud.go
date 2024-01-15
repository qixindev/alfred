package sms

import (
	"alfred/backend/model"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Tcloud struct {
	Config model.SmsTcloud
}

func (s *Tcloud) Send(number string, code string) error {
	credential := common.NewCredential(s.Config.SecretId, s.Config.SecretKey)
	client, _ := sms.NewClient(credential, s.Config.Region, profile.NewClientProfile())
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.Config.SdkAppId)
	request.SignName = common.StringPtr(s.Config.SignName)
	request.TemplateId = common.StringPtr(s.Config.TemplateId)
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{number})

	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	} else if err != nil {
		return err
	}

	for _, v := range response.Response.SendStatusSet {
		if v.Code != nil && *v.Code != "Ok" {
			return fmt.Errorf(*v.Code)
		}
	}

	return nil
}

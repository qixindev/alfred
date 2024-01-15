package sms

import (
	"alfred/backend/model"
	"fmt"
	ali "github.com/alibabacloud-go/dysmsapi-20170525/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
)

type Alibaba struct {
	Config model.SmsAlibaba
}

func (a *Alibaba) Send(phoneNumber string, code string) error {
	smsClient, err := ali.NewClient(&rpc.Config{
		AccessKeyId:     tea.String(a.Config.AccessKeyId),
		AccessKeySecret: tea.String(a.Config.AccessKeySecret),
		RegionId:        tea.String(a.Config.RegionId),
		Endpoint:        tea.String(a.Config.Endpoint),
	})
	if err != nil {
		return err
	}

	smsResp, err := smsClient.SendSms(&ali.SendSmsRequest{
		SignName:      tea.String(a.Config.SignName),
		TemplateCode:  tea.String(a.Config.TemplateCode),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
	})
	if err != nil {
		return err
	}
	println(smsResp.String())
	return nil
}

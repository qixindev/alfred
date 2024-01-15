package model

import "alfred/backend/endpoint/req"

type InterfaceSms interface {
	Save(sms req.Sms) any
}

type SmsConnector struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TenantId uint   `gorm:"primaryKey"`
	Tenant   Tenant `json:"-"`
}

type SmsTcloud struct {
	Id             uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	SmsConnectorId uint         `json:"smsId"`
	SmsConnector   SmsConnector `gorm:"foreignKey:SmsConnectorId, TenantId" json:"smsConnector"`
	SecretId       string       `json:"secretId"`
	SecretKey      string       `json:"secretKey"`
	Region         string       `json:"region"`
	SdkAppId       string       `json:"sdkAppId"`
	SignName       string       `json:"signName"`
	TemplateId     string       `json:"templateId"`

	TenantId uint `gorm:"primaryKey" json:"tenantId"`
}

func (s *SmsTcloud) Save(r req.Sms) any {
	return &SmsTcloud{
		SmsConnectorId: r.Id,
		SecretId:       r.SecretId,
		SecretKey:      r.SecretKey,
		Region:         r.Region,
		SdkAppId:       r.SdkAppId,
		SignName:       r.SignName,
		TemplateId:     r.TemplateId,
		TenantId:       r.TenantId,
	}
}

type SmsAlibaba struct {
	Id              uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	SmsConnectorId  uint         `json:"smsId"`
	SmsConnector    SmsConnector `gorm:"foreignKey:SmsConnectorId, TenantId" json:"smsConnector"`
	AccessKeyId     string       `json:"secretId"`
	AccessKeySecret string       `json:"secretKey"`
	RegionId        string       `json:"region"`
	Endpoint        string       `json:"sdkAppId"`
	SignName        string       `json:"signName"`
	TemplateCode    string       `json:"templateId"`

	TenantId uint `gorm:"primaryKey" json:"tenantId"`
}

func (s *SmsAlibaba) Save(r req.Sms) any {
	return &SmsAlibaba{
		SmsConnectorId:  r.Id,
		AccessKeyId:     r.AccessKeyId,
		AccessKeySecret: r.AccessKeySecret,
		RegionId:        r.RegionId,
		Endpoint:        r.Endpoint,
		SignName:        r.SignName,
		TemplateCode:    r.TemplateCode,
		TenantId:        r.TenantId,
	}
}

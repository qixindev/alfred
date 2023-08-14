package model

type SmsConnector struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TenantId uint   `gorm:"primaryKey"`
	Tenant   Tenant
}

type SmsTcloud struct {
	Id             uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	SmsConnectorId uint         `json:"smsId"`
	SmsConnector   SmsConnector `gorm:"foreignKey:SmsConnectorId, TenantId"`
	SecretId       string       `json:"secretId"`
	SecretKey      string       `json:"secretKey"`
	Region         string       `json:"region"`
	SdkAppId       string       `json:"sdkAppId"`
	SignName       string       `json:"signName"`
	TemplateId     string       `json:"templateId"`

	TenantId uint `gorm:"primaryKey" json:"tenantId"`
}

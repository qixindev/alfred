package models

type ProviderAzureAd struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	TenantId uint `gorm:"primaryKey"`
}

type ProviderDingTalk struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	AgentId   string
	AppKey    string
	AppSecret string

	TenantId uint `gorm:"primaryKey"`
}

type ProviderOAuth2 struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	TokenEndpoint     string `json:"tokenEndpoint"`
	UserinfoEndpoint  string `json:"userinfoEndpoint"`
	Scope             string `json:"scope"`
	ResponseType      string `json:"responseType"`

	TenantId uint `gorm:"primaryKey"`
}

type ProviderWeCom struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	CorpId    string `json:"corpId"`
	AgentId   string `json:"agentId"`
	AppSecret string `json:"appSecret"`

	TenantId uint `gorm:"primaryKey"`
}

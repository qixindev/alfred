package req

import "github.com/gin-gonic/gin"

type Provider struct {
	Id           uint   `json:"id"`
	ProviderId   uint   `json:"providerId"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	AgentId      string `json:"agentId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AppKey       string `json:"appKey"`
	AppSecret    string `json:"appSecret"`
	CorpId       string `json:"corpId"`

	AuthorizeEndpoint string `json:"authorizeEndpoint"`
	TokenEndpoint     string `json:"tokenEndpoint"`
	UserinfoEndpoint  string `json:"userinfoEndpoint"`
	Scope             string `json:"scope"`
	ResponseType      string `json:"responseType"`

	TenantId uint `json:"-"`
}

func (r *Provider) Dto() any {
	return &gin.H{
		"name": r.Name,
		"type": r.Type,
	}
}

type Sms struct {
	Id       uint   `json:"id"`
	TenantId uint   `json:"-"`
	Name     string `json:"name"`
	Type     string `json:"type"`

	SecretId   int64
	SecretKey  string `json:"secretKey"`
	Region     string `json:"region"`
	SdkAppId   string `json:"sdkAppId"`
	SignName   string `json:"signName"`
	TemplateId uint   `json:"templateId"`
}

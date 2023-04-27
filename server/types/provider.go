package types

import "github.com/gin-gonic/gin"

type ReqProvider struct {
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

func (r *ReqProvider) Dto() any {
	return &gin.H{
		"name": r.Name,
		"type": r.Type,
	}
}

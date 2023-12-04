package model

import (
	"alfred/internal/endpoint/req"
	"github.com/gin-gonic/gin"
)

type ProviderWechat struct {
	Id         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId uint     `json:"providerId"`
	Provider   Provider `gorm:"foreignKey:ProviderId, TenantId" json:"provider"`

	AppId     string `json:"agentId"`
	AppSecret string `json:"appSecret"`
	TenantId  uint   `json:"tenantId" gorm:"primaryKey"`
}

func (p *ProviderWechat) Dto() any {
	return &gin.H{
		"providerId": p.ProviderId,
		"name":       p.Provider.Name,
		"type":       p.Provider.Type,
		"agentId":    p.AppId,
		"appKey":     "",
		"appSecret":  p.AppSecret,
	}
}
func (p *ProviderWechat) Save(r req.Provider) any {
	return &ProviderWechat{
		ProviderId: r.ProviderId,
		TenantId:   r.TenantId,
		AppId:      r.AgentId,
		AppSecret:  r.AppSecret,
	}
}

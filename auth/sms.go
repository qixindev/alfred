package auth

import (
	"accounts/connectors"
	"accounts/models"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type ProviderSms struct {
	Id             uint                    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderId     uint                    `json:"providerId"`
	Provider       models.Provider         `gorm:"foreignKey:ProviderId, TenantId"`
	SmsConnectorId uint                    `json:"smsConnectorId"`
	SmsConnector   connectors.SmsConnector `gorm:"foreignKey:SmsConnectorId, TenantId"`
	TenantId       uint                    `gorm:"primaryKey;autoIncrement" json:"tenantId"`
}

type PhoneVerification struct {
	Id        uint
	Phone     string
	Code      string
	CreatedAt time.Time
}

func (p *ProviderSms) Auth(number string) string {

	return "sent"
}

func (p ProviderSms) Login(c *gin.Context) (*models.UserInfo, error) {
	return nil, errors.New("")
}

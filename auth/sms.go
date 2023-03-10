package auth

import (
	"accounts/connectors"
	"accounts/models"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type ProviderSms struct {
	Config models.ProviderSms
}

type PhoneVerification struct {
	Id        uint
	Phone     string
	Code      string
	CreatedAt time.Time
}

func (p *ProviderSms) Auth(number string) string {
	connector := connectors.GetConnector()
	return "sent"
}

func (p *ProviderSms) Login(c *gin.Context) (*models.UserInfo, error) {
	return nil, errors.New("")
}

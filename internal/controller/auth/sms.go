package auth

import (
	"accounts/internal/controller/connectors"
	"accounts/internal/model"
	"accounts/pkg/global"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ProviderSms struct {
	Config model.ProviderSms
}

func (p *ProviderSms) Auth(number string) (string, error) {
	connector, err := connectors.GetConnector(p.Config.TenantId, p.Config.SmsConnectorId)
	if err != nil {
		return "", err
	}

	var phoneVerification model.PhoneVerification
	if strings.HasPrefix(number, "+86") == false {
		return "", errors.New("only +86 suffix supported")
	}
	code := fmt.Sprint(time.Now().Nanosecond())
	if global.DB.First(&phoneVerification, "phone = ?", number).Error == nil {
		// found existing verification
		if phoneVerification.CreatedAt.Add(time.Minute).Unix() > time.Now().Unix() {
			return "", errors.New("too fast")
		}
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err := global.DB.Save(&phoneVerification).Error; err != nil {
			return "", err
		}
	} else {
		phoneVerification.Phone = number
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err := global.DB.Create(&phoneVerification).Error; err != nil {
			return "", err
		}
	}
	err = connector.Send(number, []string{code})
	if err != nil {
		return "", err
	}

	return "sent", nil
}

func (p *ProviderSms) Login(c *gin.Context) (*model.UserInfo, error) {
	phone := c.PostForm("phone")
	code := c.PostForm("code")
	var v model.PhoneVerification
	if err := global.DB.First(&v, "phone = ? AND code = ?", phone, code).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Delete(&v).Error; err != nil {
		return nil, err
	}
	u := model.UserInfo{
		Sub:   phone,
		Phone: phone,
	}
	return &u, nil
}
func (p *ProviderSms) LoginConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
	}
}

func (p *ProviderSms) ProviderConfig() *gin.H {
	return &gin.H{
		"providerId": p.Config.ProviderId,
	}
}

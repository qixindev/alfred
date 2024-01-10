package auth

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ProviderSms struct {
	Config model.ProviderSms
}

func (p *ProviderSms) Auth(_ string, number string, _ uint) (string, error) {
	connector, err := service.GetConnector(p.Config.TenantId, p.Config.SmsConnectorId)
	if err != nil {
		return "", err
	}

	var phoneVerification model.PhoneVerification
	if strings.HasPrefix(number, "+86") == false {
		return "", errors.New("only +86 suffix supported")
	}
	code := utils.GetCode()
	if global.DB.First(&phoneVerification, "phone = ?", number).Error == nil {
		// found existing verification
		if phoneVerification.CreatedAt.Add(time.Minute).Unix() > time.Now().Unix() {
			return "", errors.New("too fast")
		}
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err = global.DB.Save(&phoneVerification).Error; err != nil {
			return "", err
		}
	} else {
		phoneVerification.Phone = number
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err = global.DB.Create(&phoneVerification).Error; err != nil {
			return "", err
		}
	}
	if err = connector.Send(number, []string{code}); err != nil {
		return "", err
	}

	return "", nil
}

func (p *ProviderSms) Login(c *gin.Context) (*model.UserInfo, error) {
	phone := c.Query("phone")
	code := c.Query("code")
	if phone == "" || code == "" {
		return nil, errors.New("invalid phone or code")
	}
	var v model.PhoneVerification
	if err := global.DB.First(&v, "phone = ? AND code = ?", phone, code).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Delete(&v).Error; err != nil {
		return nil, err
	}
	phoneNumber := phone
	if strings.HasPrefix(phoneNumber, "+86") {
		phoneNumber = phoneNumber[3:]
	}
	u := model.UserInfo{
		Name:        phoneNumber,
		Sub:         phone,
		Phone:       phoneNumber,
		DisplayName: "-",
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

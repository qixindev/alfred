package reset

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service/sms"
	"errors"
	"strings"
	"time"
)

type ProviderResetSms struct {
	Config model.ProviderSms
}

func (p *ProviderResetSms) ResetAuth(number string, _ uint) (location string, code string, error error) {
	connector, err := sms.GetConnector(p.Config.TenantId, p.Config.SmsConnectorId)
	if err != nil {
		return "", "", err
	}

	var phoneVerification model.PhoneVerification
	if strings.HasPrefix(number, "+86") == false {
		return "", "", errors.New("only +86 suffix supported")
	}
	code = utils.GetCode()
	if global.DB.First(&phoneVerification, "phone = ?", number).Error == nil {
		// found existing verification
		if phoneVerification.CreatedAt.Add(time.Minute).Unix() > time.Now().Unix() {
			return "", "", errors.New("too fast")
		}
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err = global.DB.Save(&phoneVerification).Error; err != nil {
			return "", "", err
		}

	} else {
		phoneVerification.Phone = number
		phoneVerification.Code = code
		phoneVerification.CreatedAt = time.Now()
		if err = global.DB.Create(&phoneVerification).Error; err != nil {
			return "", "", err
		}
	}
	if err = connector.Send(number, code); err != nil {
		return "", "", err
	}
	return "", code, nil
}

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

	if strings.HasPrefix(number, "+86") == false {
		return "", "", errors.New("only +86 suffix supported")
	}
	code = utils.GetCode()
	oldCode, err := global.GetCodeCache(number)
	if err == nil {
		if oldCode.Time.Add(time.Minute).Unix() > time.Now().Unix() {
			return "", "", errors.New("code send too fast")
		}
	}
	if err = connector.Send(number, code); err != nil {
		return "", "", err
	}

	if err = global.SetCodeCache(number, code); err != nil {
		return "", "", err
	}
	return "", code, nil
}

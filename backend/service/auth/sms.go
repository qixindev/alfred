package auth

import (
	"alfred/backend/model"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"alfred/backend/service/sms"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ProviderSms struct {
	Config model.ProviderSms
}

func (p *ProviderSms) Auth(_ string, phone string, _ uint) (string, error) {
	connector, err := sms.GetConnector(p.Config.TenantId, p.Config.SmsConnectorId)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(phone, "+86") == false {
		return "", errors.New("only +86 suffix supported")
	}
	code := utils.GetCode()
	oldCode, err := global.GetCodeCache(phone)
	if err == nil {
		if oldCode.Time.Add(time.Minute).Unix() > time.Now().Unix() {
			return "", errors.New("code send too fast")
		}
	}
	if err = connector.Send(phone, code); err != nil {
		return "", err
	}
	if err = global.SetCodeCache(phone, code); err != nil {
		return "", err
	}

	return code, nil
}

func (p *ProviderSms) Login(code string, loginInfo global.StateInfo) (*model.UserInfo, error) {
	phone := loginInfo.AuthState
	codeInfo, err := global.GetCodeCache(phone)
	if err != nil {
		return nil, err
	}
	if code != codeInfo.Code {
		return nil, errors.New("invalid login code")
	}
	if err = global.CodeCache.Delete(phone); err != nil {
		return nil, err
	}

	return &model.UserInfo{
		Name:        phone,
		Sub:         phone,
		Phone:       phone,
		DisplayName: "-",
	}, nil
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

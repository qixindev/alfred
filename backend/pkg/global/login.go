package global

import (
	"encoding/json"
	"time"
)

type StateInfo struct {
	State      string `json:"state"`
	AuthState  string `json:"authState"`  // phon or email
	AuthString string `json:"authString"` // location or code
	Type       string `json:"type"`       // ding, wecom, sms
	Provider   string `json:"provider"`   // provider name
	ClientId   string `json:"clientId"`   // client
	Tenant     string `json:"tenant"`     // tenant name
	TenantId   uint   `json:"tenantId"`   // tenant id
}

func getStateInfo(state string) (StateInfo, error) {
	loginInfo, err := StateCache.Get(state)
	if err != nil {
		return StateInfo{}, err
	}
	var stateInfo StateInfo
	if err = json.Unmarshal(loginInfo, &stateInfo); err != nil {
		return StateInfo{}, err
	}
	return stateInfo, nil
}

func GetAndDeleteStateInfo(state string) (StateInfo, error) {
	loginInfo, err := getStateInfo(state)
	if err != nil {
		return StateInfo{}, err
	}
	return loginInfo, StateCache.Delete(state)
}

func SetStateInfo(state string, loginInfo StateInfo) error {
	infoByte, err := json.Marshal(&loginInfo)
	if err != nil {
		return err
	}
	if err = StateCache.Set(state, infoByte); err != nil {
		return err
	}
	return nil
}

type CodeInfo struct {
	Object string    `json:"object"` // phone or email
	Code   string    `json:"code"`   // code
	Time   time.Time `json:"time"`   // created timestamp
}

func GetCodeCache(obj string) (CodeInfo, error) {
	code, err := CodeCache.Get(obj)
	if err != nil {
		return CodeInfo{}, err
	}
	var stateInfo CodeInfo
	if err = json.Unmarshal(code, &stateInfo); err != nil {
		return CodeInfo{}, err
	}
	return stateInfo, nil
}

func SetCodeCache(obj string, code string) error {
	infoByte, err := json.Marshal(&CodeInfo{Object: obj, Code: code, Time: time.Now()})
	if err != nil {
		return err
	}
	return CodeCache.Set(obj, infoByte)
}

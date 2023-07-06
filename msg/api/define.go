package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Wecom struct {
	AgentId int64  `mapstructure:"agentId" json:"agentId" yaml:"agentId"`
	CorpId  string `mapstructure:"corpId" json:"corpId" yaml:"corpId"`
	Secret  string `mapstructure:"secret" json:"secret" yaml:"secret"`
}
type Ding struct {
	AgentId   int64  `mapstructure:"agentId" json:"agentId" yaml:"agentId"`
	AppKey    string `mapstructure:"appKey" json:"appKey" yaml:"appKey"`
	AppSecret string `mapstructure:"appSecret" json:"appSecret" yaml:"appSecret"`
}
type Third struct {
	Login struct {
		LoginType string `mapstructure:"loginType" json:"loginType" yaml:"loginType"`
		Host      string `mapstructure:"host" json:"host" yaml:"host"`
	} `mapstructure:"login" json:"login" yaml:"login"`
	Wecom Wecom `mapstructure:"wecom" json:"wecom" yaml:"wecom"`
	Ding  Ding  `mapstructure:"ding" json:"ding" yaml:"ding"`
}

func toInt64(s any) int64 {
	res, ok := s.(string)
	if !ok {
		return 0
	}
	parseInt, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return 0
	}

	return parseInt
}

func GetWecomConfig(c gin.H) (*Wecom, error) {
	res := Wecom{
		AgentId: toInt64(c["agentId"]),
		CorpId:  c["corpId"].(string),
		Secret:  c["secret"].(string),
	}

	if res.AgentId == 0 || res.CorpId == "" || res.Secret == "" {
		return nil, errors.New("invalid wecom config")
	}

	return &res, nil
}

func GetDingTalkConfig(c gin.H) (*Ding, error) {
	res := Ding{
		AgentId:   toInt64(c["agentId"]),
		AppKey:    c["appKey"].(string),
		AppSecret: c["appSecret"].(string),
	}

	if res.AgentId == 0 || res.AppKey == "" || res.AppSecret == "" {
		return nil, errors.New("invalid ding talk config")
	}

	return &res, nil
}

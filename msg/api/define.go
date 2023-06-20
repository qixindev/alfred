package api

type Wecom struct {
	AgentId int    `mapstructure:"agentId" json:"agentId" yaml:"agentId"`
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

// ==========企业微信配置==========

func (t *Third) GetWecomAgentId() int {
	return t.Wecom.AgentId
}
func (t *Third) GetWecomCorpId() string {
	return t.Wecom.CorpId
}
func (t *Third) GetWecomSecret() string {
	return t.Wecom.Secret
}

// ==========钉钉配置==========

func (t *Third) GetDingAppKey() string {
	return t.Ding.AppKey
}
func (t *Third) GetDingAppSecret() string {
	return t.Ding.AppSecret
}
func (t *Third) GetDingAgentId() int64 {
	return t.Ding.AgentId
}

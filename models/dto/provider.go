package dto

type ProviderDto struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ProviderConfigDto struct {
	ProviderId   uint   `json:"providerId"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	AgentId      string `json:"agentId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

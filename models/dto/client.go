package dto

type ClientDto struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
}

type RedirectUriDto struct {
	Id          uint   `json:"id"`
	RedirectUri string `json:"redirectUri"`
}

type ClientSecretDto struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type AccessTokenDto struct {
	AccessToken string `json:"access_token"`
}

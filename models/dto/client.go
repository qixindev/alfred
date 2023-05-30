package dto

type ClientDto struct {
	Id       string `json:"id"`
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

type ClientUserDto struct {
	Id       uint   `json:"id"`
	Sub      string `json:"sub"`
	ClientId string `json:"clientId"`
	UserName string `json:"userName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type AccessTokenDto struct {
	AccessToken string `json:"access_token"`
}

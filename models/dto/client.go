package dto

type ClientDto struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
}

type RedirectUriDto struct {
	Id          uint   `gorm:"json:"id"`
	RedirectUri string `json:"redirectUri"`
}

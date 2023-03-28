package dto

type DeviceDto struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type DeviceSecretDto struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

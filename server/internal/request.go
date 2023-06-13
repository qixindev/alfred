package internal

type IamNameRequest struct {
	Name string `json:"name"`
}

type IamActionRequest struct {
	ActionId string `json:"actionId"`
}

type IamUserRequest struct {
	UserId uint `json:"userId"`
}

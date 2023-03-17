package dto

type GroupMemberDto struct {
	Type string `json:"type,omitempty"`
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role,omitempty"`
}

type GroupDto struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	ParentId uint   `json:"parentId"`
}

type GroupDeviceDto struct {
	Id       uint `json:"id"`
	GroupId  uint `json:"name"`
	DeviceId uint `json:"deviceId"`
}

type GroupUserDto struct {
	Id      uint   `json:"id"`
	GroupId uint   `json:"groupId"`
	UserId  uint   `json:"userId"`
	Role    string `json:"role"`
}

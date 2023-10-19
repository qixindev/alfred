package req

type IamResourceId struct {
	Tenant   string `json:"tenant" uri:"tenant"`
	Client   string `json:"client" uri:"client"`
	TypeId   string `json:"typeId" uri:"typeId"`
	RoleId   string `json:"roleId" uri:"roleId"`
	ActionId string `json:"actionId" uri:"actionId"`
}

type IamClientUser struct {
	Tenant       string `json:"tenant" uri:"tenant"`
	Client       string `json:"client" uri:"client"`
	TypeId       string `json:"typeId" uri:"typeId"`
	RoleId       string `json:"roleId" uri:"roleId"`
	ActionId     string `json:"actionId" uri:"actionId"`
	ClientUserId []struct {
		UserId uint `json:"userId"`
	} `json:"clientUserId"`
}

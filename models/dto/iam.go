package dto

type ResourceTypeRoleActionDto struct {
	Id         uint   `json:"id"`
	RoleId     uint   `json:"roleId"`
	TenantId   uint   `json:"tenantId"`
	ActionId   uint   `json:"actionId"`
	ActionName string `json:"actionName"`
}

type ResourceRoleUserDto struct {
	ResourceName string `json:"resourceName"`
	Role         string `json:"role"`
	Sub          string `json:"sub"`
}

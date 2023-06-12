package dto

type ResourceTypeRoleActionDto struct {
	Id         uint   `json:"id"`
	RoleId     string `json:"roleId"`
	TenantId   uint   `json:"tenantId"`
	ActionId   string `json:"actionId"`
	ActionName string `json:"actionName"`
}

type ResourceRoleUserDto struct {
	Id           uint   `json:"id"`
	ResourceName string `json:"resourceName"`
	RoleName     string `json:"roleName"`
	DisplayName  string `json:"userName"`
	Sub          string `json:"sub"`
}

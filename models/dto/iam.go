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
	ResourceId   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	RoleId       string `json:"roleId"`
	RoleName     string `json:"roleName"`
	DisplayName  string `json:"displayName"`
	Sub          string `json:"sub"`
}

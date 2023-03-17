package dto

type ResourceTypeRoleActionDto struct {
	Id         uint   `json:"id"`
	RoleId     uint   `json:"roleId"`
	TenantId   uint   `json:"tenantId"`
	ActionId   uint   `json:"actionId"`
	ActionName string `json:"actionName"`
}

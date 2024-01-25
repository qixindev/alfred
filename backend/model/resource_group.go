package model

type (
	ResourceGroup struct {
		Id          string `gorm:"primaryKey" json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ClientId    string `json:"clientId" uri:"client"`
		Client      Client `gorm:"foreignKey:ClientId, TenantId" json:"-" swaggerignore:"true"`
		TenantId    uint   `gorm:"primaryKey"`
	}

	ResourceGroupResource struct {
		Id          string        `gorm:"primaryKey" json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		GroupId     string        `json:"groupId" uri:"groupId"`
		Group       ResourceGroup `gorm:"foreignKey:GroupId, TenantId" json:"-" swaggerignore:"true"`
		TenantId    uint          `gorm:"primaryKey"`
	}
	ResourceGroupRole struct {
		Id          string        `gorm:"primaryKey" json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		GroupId     string        `json:"groupId" uri:"groupId"`
		Group       ResourceGroup `gorm:"foreignKey:GroupId, TenantId" json:"-" swaggerignore:"true"`
		TenantId    uint          `gorm:"primaryKey"`
	}
	ResourceGroupAction struct {
		Id          string        `gorm:"primaryKey" json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		GroupId     string        `json:"groupId" uri:"groupId"`
		Group       ResourceGroup `gorm:"foreignKey:GroupId, TenantId" json:"-" swaggerignore:"true"`
		TenantId    uint          `gorm:"primaryKey"`
	}

	ResourceGroupRoleAction struct {
		Id         uint                `gorm:"primaryKey" json:"id"`
		RoleId     string              `json:"roleId" uri:"roleId"`
		Role       ResourceGroupRole   `gorm:"foreignKey:RoleId, TenantId" json:"-" swaggerignore:"true"`
		ActionId   string              `json:"actionId" uri:"actionId"`
		Action     ResourceGroupAction `gorm:"foreignKey:ActionId, TenantId" json:"-" swaggerignore:"true"`
		ActionName string              `gorm:"<-:false;-:migration" json:"actionName"`
		RoleName   string              `gorm:"<-:false;-:migration" json:"roleName"`
		TenantId   uint                `gorm:"primaryKey"`
	}
	ResourceGroupUser struct {
		Id                uint              `gorm:"primaryKey" json:"id"`
		GroupId           string            `json:"groupId" uri:"groupId"`
		ResourceGroup     ResourceGroup     `gorm:"foreignKey:GroupId, TenantId" json:"-" swaggerignore:"true"`
		ResourceGroupName string            `json:"resourceGroupName" gorm:"<-:false;-:migration"`
		RoleId            string            `json:"roleId" uri:"roleId"`
		Role              ResourceGroupRole `gorm:"foreignKey:RoleId, TenantId" json:"role" swaggerignore:"true"`
		RoleName          string            `json:"roleName" gorm:"<-:false;-:migration"`
		ClientUserId      uint              `json:"userId" uri:"userId"`
		ClientUser        ClientUser        `gorm:"foreignKey:ClientUserId, TenantId" json:"user" swaggerignore:"true"`
		TenantId          uint              `gorm:"primaryKey"`
		Sub               string            `json:"sub" gorm:"<-:false;-:migration"`
		DisplayName       string            `json:"displayName" gorm:"<-:false;-:migration"`
	}
)

type RequestResourceGroup struct {
	Tenant      Tenant   `json:"-" swaggerignore:"true"`
	ClientId    string   `json:"clientId" uri:"client" swaggerignore:"true"`
	GroupId     string   `json:"groupId" uri:"groupId" swaggerignore:"true"`
	ResourceId  string   `json:"resourceId" uri:"resourceId" swaggerignore:"true"`
	ActionId    string   `json:"actionId" uri:"actionId" swaggerignore:"true"`
	RoleId      string   `json:"roleId" uri:"roleId"`
	UserId      uint     `json:"userId" uri:"userId" swaggerignore:"true"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Uid         string   `json:"uid"`
	ActionIds   []string `json:"actionIds"`
}

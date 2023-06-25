package iam

import (
	"accounts/models"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
)

func GetClientFromCid(c *gin.Context) (*models.Client, error) {
	cid := c.Param("client")
	var client models.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", cid).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func AddIamRoutes(rg *gin.RouterGroup) {
	rg.GET("/types", ListIamType)
	rg.POST("/types", NewIamType)
	rg.DELETE("/types/:typeId", DeleteIamType)

	// 资源管理
	rg.GET("/types/:typeId/resources", ListIamResource)
	rg.POST("/types/:typeId/resources", NewIamResource)
	rg.DELETE("/types/:typeId/resources/:resourceId", DeleteIamResource)

	// 角色管理
	rg.GET("/types/:typeId/roles", ListIamRole)
	rg.POST("/types/:typeId/roles", NewIamRole)
	rg.DELETE("/types/:typeId/roles/:roleId", DeleteIamRole)

	// 动作管理
	rg.GET("/types/:typeId/actions", ListIamAction)
	rg.POST("/types/:typeId/actions", NewIamAction)
	rg.DELETE("/types/:typeId/actions/:actionId", DeleteIamAction)

	// 针对某类资源，授予角色动作
	rg.GET("/types/:typeId/roles/:roleId/actions", ListIamRoleAction)
	rg.POST("/types/:typeId/roles/:roleId/actions", NewIamRoleAction)
	rg.DELETE("/types/:typeId/roles/:roleId/actions/:actionId", DeleteIamRoleAction)

	// 针对某个资源，对用户赋予角色
	rg.GET("/types/:typeId/resources/:resourceId/roles/:roleId/users", ListIamResourceRole)
	rg.POST("/types/:typeId/resources/:resourceId/roles/:roleId/users", NewIamResourceRole)
	rg.DELETE("/types/:typeId/resources/:resourceId/roles/:roleId/users/:userId", DeleteIamResourceRoleUser)

	rg.GET("/types/:typeId/actions/:actionId/users/:user/resources", GetIamActionResource) // 针对一类资源，用户拥有哪些资源的哪些角色
	rg.GET("/types/:typeId/resources/:resourceId/actions/:actionId/users/:user", IsUserActionPermission)
}

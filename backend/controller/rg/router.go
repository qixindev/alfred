package rg

import "github.com/gin-gonic/gin"

func AddResourceGroupRoutes(rg *gin.RouterGroup) {
	// 资源组管理
	rg.GET("/resourceGroups", GetResourceGroupList)
	rg.POST("/resourceGroups", CreateResourceGroup)
	rg.GET("/resourceGroups/:groupId", GetResourceGroup)
	rg.PUT("/resourceGroups/:groupId", UpdateResourceGroup)
	rg.DELETE("/resourceGroups/:groupId", DeleteResourceGroup)

	// 组内资源管理
	rg.GET("/resourceGroups/:groupId/resources", GetResourceGroupResourceList)
	rg.POST("/resourceGroups/:groupId/resources", CreateResourceGroupResource)
	rg.GET("/resourceGroups/:groupId/resources/:resourceId", GetResourceGroupResource)
	rg.PUT("/resourceGroups/:groupId/resources/:resourceId", UpdateResourceGroupResource)
	rg.DELETE("/resourceGroups/:groupId/resources/:resourceId", DeleteResourceGroupResource)

	// 资源组角色管理
	rg.GET("/resourceGroups/:groupId/roles", GetResourceGroupRoleList)
	rg.POST("/resourceGroups/:groupId/roles", CreateResourceGroupRole)
	rg.GET("/resourceGroups/:groupId/roles/:roleId", GetResourceGroupRole)
	rg.PUT("/resourceGroups/:groupId/roles/:roleId", UpdateResourceGroupRole)
	rg.DELETE("/resourceGroups/:groupId/roles/:roleId", DeleteResourceGroupRole)

	// 资源组动作管理
	rg.GET("/resourceGroups/:groupId/actions", GetResourceGroupActionList)
	rg.POST("/resourceGroups/:groupId/actions", CreateResourceGroupAction)
	rg.GET("/resourceGroups/:groupId/actions/:actionId", GetResourceGroupAction)
	rg.PUT("/resourceGroups/:groupId/actions/:actionId", UpdateResourceGroupAction)
	rg.DELETE("/resourceGroups/:groupId/actions/:actionId", DeleteResourceGroupAction)

	// 角色动作管理
	rg.GET("/resourceGroups/:groupId/roles/:roleId/actions", GetResourceGroupRoleActionList)       // 角色所拥有的权限列表
	rg.POST("/resourceGroups/:groupId/roles/:roleId/actions", CreateResourceGroupRoleAction)       // 给角色新增权限，支持数组
	rg.GET("/resourceGroups/:groupId/roles/:roleId/actions/:actionId", GetResourceGroupRoleAction) // 角色是否拥有某个权限
	rg.PUT("/resourceGroups/:groupId/roles/:roleId/actions", UpdateResourceGroupRoleAction)        // 修改角色的权限，支持数组
	rg.DELETE("/resourceGroups/:groupId/roles/:roleId/actions", DeleteResourceGroupRoleAction)     // 删除角色的权限，支持数组

	// 用户在组内角色
	rg.GET("/resourceGroups/:groupId/users/:userId/role", GetResourceGroupUserRole)                // 用户在组内的角色
	rg.GET("/resourceGroups/:groupId/users/:userId/actions", GetResourceGroupUserActionList)       // 用户在组内所拥有的权限列表
	rg.GET("/resourceGroups/:groupId/users/:userId/actions/:actionId", GetResourceGroupUserAction) // 用户在组内是否拥有某个权限
	rg.POST("/resourceGroups/:groupId/users/:userId", CreateResourceGroupUserRole)                 // 将用户拉入组内
	rg.PUT("/resourceGroups/:groupId/users/:userId", UpdateResourceGroupUserRole)                  // 修改用户在组内的角色
	rg.DELETE("/resourceGroups/:groupId/users/:userId", DeleteResourceGroupUser)                   // 踢出用户
}

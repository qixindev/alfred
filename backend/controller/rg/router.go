package rg

import "github.com/gin-gonic/gin"

func AddResourceGroupRoutes(rg *gin.RouterGroup) {
	// 资源组管理
	rg.GET("/resourceGroups", CreateResourceGroup)
	rg.POST("/resourceGroups", CreateResourceGroup)
	rg.GET("/resourceGroups/:groupId", CreateResourceGroup)
	rg.PUT("/resourceGroups/:groupId", CreateResourceGroup)
	rg.DELETE("/resourceGroups/:groupId", CreateResourceGroup)

	// 组内资源管理
	rg.GET("/resourceGroups/:groupId/resources", CreateResourceGroup)
	rg.POST("/resourceGroups/:groupId/resources", CreateResourceGroup)
	rg.GET("/resourceGroups/:groupId/resources/:resourceId", CreateResourceGroup)
	rg.PUT("/resourceGroups/:groupId/resources/:resourceId", CreateResourceGroup)
	rg.DELETE("/resourceGroups/:groupId/resources/:resourceId", CreateResourceGroup)

	// 资源组角色管理
	rg.GET("/resourceGroups/:groupId/roles", CreateResourceGroup)
	rg.POST("/resourceGroups/:groupId/roles", CreateResourceGroup)
	rg.GET("/resourceGroups/:groupId/roles/:roleId", CreateResourceGroup)
	rg.PUT("/resourceGroups/:groupId/roles/:roleId", CreateResourceGroup)
	rg.DELETE("/resourceGroups/:groupId/roles/:roleId", CreateResourceGroup)

	// 角色动作管理
	rg.GET("/resourceGroups/:groupId/roles/:roleId/actions", CreateResourceGroup)           // 角色所拥有的权限列表
	rg.POST("/resourceGroups/:groupId/roles/:roleId/actions", CreateResourceGroup)          // 给角色新增权限，支持数组
	rg.GET("/resourceGroups/:groupId/roles/:roleId/actions/:actionId", CreateResourceGroup) // 角色是否拥有某个权限
	rg.PUT("/resourceGroups/:groupId/roles/:roleId/actions", CreateResourceGroup)           // 修改角色的权限，支持数组
	rg.DELETE("/resourceGroups/:groupId/roles/:roleId/actions", CreateResourceGroup)        // 删除角色的权限，支持数组

	// 用户在组内角色
	rg.GET("/resourceGroups/:groupId/users/:userId/roles", CreateResourceGroup)             // 用户在组内的角色
	rg.GET("/resourceGroups/:groupId/users/:userId/actions", CreateResourceGroup)           // 用户在组内所拥有的权限列表
	rg.GET("/resourceGroups/:groupId/users/:userId/actions/:actionId", CreateResourceGroup) // 用户在组内是否拥有某个权限
	rg.POST("/resourceGroups/:groupId/users/:userId", CreateResourceGroup)                  // 将用户拉入组内
	rg.PUT("/resourceGroups/:groupId/users/:userId", CreateResourceGroup)                   // 修改用户在组内的角色
	rg.DELETE("/resourceGroups/:groupId/users/:userId", CreateResourceGroup)                // 踢出用户

}

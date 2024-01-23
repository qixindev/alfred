package rg

import "github.com/gin-gonic/gin"

func AddResourceGroupRoutes(r *gin.RouterGroup) {
	rg := r.Group("/resourceGroups")
	{
		// 资源组管理
		rg.GET("", CreateResourceGroup)
		rg.POST("/", CreateResourceGroup)
		rg.GET("/:groupId", CreateResourceGroup)
		rg.PUT("/:groupId", CreateResourceGroup)
		rg.DELETE("/:groupId", CreateResourceGroup)

		// 组内资源管理
		rg.GET("/:groupId/resources", CreateResourceGroup)
		rg.POST("/:groupId/resources", CreateResourceGroup)
		rg.GET("/:groupId/resources/:resourceId", CreateResourceGroup)
		rg.PUT("/:groupId/resources/:resourceId", CreateResourceGroup)
		rg.DELETE("/:groupId/resources/:resourceId", CreateResourceGroup)

		// 资源组角色管理
		rg.GET("/:groupId/roles", CreateResourceGroup)
		rg.POST("/:groupId/roles", CreateResourceGroup)
		rg.GET("/:groupId/roles/:roleId", CreateResourceGroup)
		rg.PUT("/:groupId/roles/:roleId", CreateResourceGroup)
		rg.DELETE("/:groupId/roles/:roleId", CreateResourceGroup)

		// 角色动作管理
		rg.GET("/:groupId/roles/:roleId/actions", CreateResourceGroup)           // 角色所拥有的权限列表
		rg.POST("/:groupId/roles/:roleId/actions", CreateResourceGroup)          // 给角色新增权限，支持数组
		rg.GET("/:groupId/roles/:roleId/actions/:actionId", CreateResourceGroup) // 角色是否拥有某个权限
		rg.PUT("/:groupId/roles/:roleId/actions", CreateResourceGroup)           // 修改角色的权限，支持数组
		rg.DELETE("/:groupId/roles/:roleId/actions", CreateResourceGroup)        // 删除角色的权限，支持数组

		// 用户在组内角色
		rg.GET("/:groupId/users/:userId/roles", CreateResourceGroup)             // 用户在组内的角色
		rg.GET("/:groupId/users/:userId/actions", CreateResourceGroup)           // 用户在组内所拥有的权限列表
		rg.GET("/:groupId/users/:userId/actions/:actionId", CreateResourceGroup) // 用户在组内是否拥有某个权限
		rg.POST("/:groupId/users/:userId", CreateResourceGroup)                  // 将用户拉入组内
		rg.PUT("/:groupId/users/:userId", CreateResourceGroup)                   // 修改用户在组内的角色
		rg.DELETE("/:groupId/users/:userId", CreateResourceGroup)                // 踢出用户
	}
}

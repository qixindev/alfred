package iam

import (
	"accounts/models"
	"accounts/server/internal"
	"github.com/gin-gonic/gin"
)

func GetClientFromCid(c *gin.Context) (*models.Client, error) {
	cid := c.Param("client")
	var client models.Client
	if err := internal.TenantDB(c).First(&client, "cli_id = ?", cid).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func getType(c *gin.Context) (*models.ResourceType, error) {
	client, err := GetClientFromCid(c)
	if err != nil {
		return nil, err
	}
	typeName := c.Param("type")
	var typ models.ResourceType
	if err := internal.TenantDB(c).First(&typ, "client_id = ? AND name = ?", client.Id, typeName).Error; err != nil {
		return nil, err
	}
	return &typ, nil
}

func getRole(c *gin.Context) (*models.ResourceTypeRole, error) {
	typ, err := getType(c)
	if err != nil {
		return nil, err
	}
	roleName := c.Param("role")
	var role models.ResourceTypeRole
	if err := internal.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func getAction(c *gin.Context) (*models.ResourceTypeAction, error) {
	typ, err := getType(c)
	if err != nil {
		return nil, err
	}
	actionName := c.Param("action")
	var action models.ResourceTypeAction
	if err := internal.TenantDB(c).First(&action, "type_id = ? AND name = ?", typ.Id, actionName).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func getResource(c *gin.Context) (*models.Resource, error) {
	typ, err := getType(c)
	if err != nil {
		return nil, err
	}
	resourceName := c.Param("resource")
	var resource models.Resource
	if err := internal.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func getResourceAndRole(c *gin.Context) (*models.Resource, *models.ResourceTypeRole, error) {
	typ, err := getType(c)
	if err != nil {
		return nil, nil, err
	}
	resourceName := c.Param("resource")
	var resource models.Resource
	if err := internal.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
		return nil, nil, err
	}

	roleName := c.Param("role")
	var role models.ResourceTypeRole
	if err := internal.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
		return nil, nil, err
	}
	return &resource, &role, nil
}

func AddIamRoutes(rg *gin.RouterGroup) {
	rg.GET("/types", ListIamResourceType)
	rg.POST("/types", NewIamResourceType)
	rg.DELETE("/types/:type", DeleteIamResourceType)

	rg.GET("/types/:type/roles", ListIamRole)
	rg.POST("/types/:type/roles", NewIamRole)
	rg.DELETE("/types/:type/roles/:role", DeleteIamRole)

	rg.GET("/types/:type/actions", ListIamAction)
	rg.POST("/types/:type/actions", NewIamAction)
	rg.DELETE("/types/:type/actions/:action", DeleteIamAction)

	rg.GET("/types/:type/roles/:role/actions", ListIamRoleAction)
	rg.POST("/types/:type/roles/:role/actions", NewIamRoleAction)
	rg.DELETE("/types/:type/roles/:role/actions/:action", DeleteIamRoleAction)

	rg.GET("/types/:type/resources", ListIamResource)
	rg.POST("/types/:type/resources", NewIamResource)
	rg.DELETE("/types/:type/resources/:resource", DeleteIamResource)

	rg.GET("/types/:type/resources/:resource/roles/:role/users", ListIamResourceRole)
	rg.POST("/types/:type/resources/:resource/roles/:role/users", NewIamResourceRole)
	rg.DELETE("/types/:type/resources/:resource/roles/:role/users/:user", DeleteIamResourceRole)

	rg.GET("/types/:type/resources/:resource/actions/:action/users/:user", GetIamActionUser)
}

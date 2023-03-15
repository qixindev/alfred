package iam

import (
	"accounts/middlewares"
	"accounts/models"
	"github.com/gin-gonic/gin"
)

func GetClientFromCid(c *gin.Context) (*models.Client, error) {
	cid := c.Param("client")
	var client models.Client
	if err := middlewares.TenantDB(c).First(&client, "cid = ?", cid).Error; err != nil {
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
	if err := middlewares.TenantDB(c).First(&typ, "client_id = ? AND name = ?", client.Id, typeName).Error; err != nil {
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
	if err := middlewares.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
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
	if err := middlewares.TenantDB(c).First(&action, "type_id = ? AND name = ?", typ.Id, actionName).Error; err != nil {
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
	if err := middlewares.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
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
	if err := middlewares.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
		return nil, nil, err
	}

	roleName := c.Param("role")
	var role models.ResourceTypeRole
	if err := middlewares.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
		return nil, nil, err
	}
	return &resource, &role, nil
}

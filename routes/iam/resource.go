package iam

import (
	"accounts/iam"
	"accounts/middlewares"
	"accounts/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

func AddIamRoutes(rg *gin.RouterGroup) {
	rg.GET("/types", func(c *gin.Context) {
		client, err := GetClientFromCid(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		types, err := iam.ListResourceTypes(client.TenantId, client.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, types)
	})
	rg.POST("/types", func(c *gin.Context) {
		var typ models.ResourceType
		if c.BindJSON(&typ) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		client, err := GetClientFromCid(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		t, err := iam.CreateResourceType(client.TenantId, client.Id, &typ)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, t)
	})
	rg.DELETE("/types/:type", func(c *gin.Context) {
		client, err := GetClientFromCid(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		typeName := c.Param("type")
		var typ models.ResourceType
		if err := middlewares.TenantDB(c).First(&typ, "client_id = ? AND name = ?", client.Id, typeName).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if err := iam.DeleteResourceType(client.TenantId, typ.Id); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/roles", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		roles, err := iam.ListResourceTypeRoles(typ.TenantId, typ.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, roles)
	})
	rg.POST("/types/:type/roles", func(c *gin.Context) {
		var role models.ResourceTypeRole
		if c.BindJSON(&role) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		r, err := iam.CreateResourceTypeRole(typ.TenantId, typ.Id, &role)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, r)
	})
	rg.DELETE("/types/:type/roles/:role", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		roleName := c.Param("role")
		var role models.ResourceTypeRole
		if err := middlewares.TenantDB(c).First(&role, "type_id = ? AND name = ?", typ.Id, roleName).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if err := iam.DeleteResourceTypeRole(typ.TenantId, role.Id); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/actions", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		actions, err := iam.ListResourceTypeActions(typ.TenantId, typ.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, actions)
	})
	rg.POST("/types/:type/actions", func(c *gin.Context) {
		var action models.ResourceTypeAction
		if c.BindJSON(&action) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		a, err := iam.CreateResourceTypeAction(typ.TenantId, typ.Id, &action)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, a)
	})
	rg.DELETE("/types/:type/actions/:action", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		actionName := c.Param("action")
		var action models.ResourceTypeAction
		if err := middlewares.TenantDB(c).First(&action, "type_id = ? AND name = ?", typ.Id, actionName).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if err := iam.DeleteResourceTypeAction(typ.TenantId, action.Id); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/roles/:role/actions", func(c *gin.Context) {
		role, err := getRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		roleActions, err := iam.ListResourceTypeRoleActions(role.TenantId, role.TypeId, role.Id)
		c.JSON(http.StatusOK, roleActions)
	})
	rg.POST("/types/:type/roles/:role/actions", func(c *gin.Context) {
		role, err := getRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		var roleAction models.ResourceTypeRoleAction
		if c.BindJSON(&roleAction) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		ra, err := iam.CreateResourceTypeRoleAction(role.TenantId, role.Id, &roleAction)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ra)
	})
	rg.DELETE("/types/:type/roles/:role/actions/:action", func(c *gin.Context) {
		role, err := getRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		actionName := c.Param("action")
		var action models.ResourceTypeAction
		if err := middlewares.TenantDB(c).First(&action, "type_id = ? AND name = ?", role.TypeId, actionName).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		var roleAction models.ResourceTypeRoleAction
		if err := middlewares.TenantDB(c).First(&roleAction, "role_id = ? AND action_id = ?", role.Id, action.Id).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if iam.DeleteResourceTypeRoleAction(role.TenantId, roleAction.Id) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/resources", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		resources, err := iam.ListResources(typ.TenantId, typ.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, resources)
	})
	rg.POST("/types/:type/resources", func(c *gin.Context) {
		var resource models.Resource
		if c.BindJSON(&resource) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		r, err := iam.CreateResource(typ.TenantId, typ.Id, &resource)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, r)
	})
	rg.DELETE("/types/:type/resources/:resource", func(c *gin.Context) {
		typ, err := getType(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		resourceName := c.Param("resource")
		var resource models.Resource
		if err := middlewares.TenantDB(c).First(&resource, "type_id = ? AND name = ?", typ.Id, resourceName).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if err := iam.DeleteResource(typ.TenantId, resource.Id); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/resources/:resource/roles/:role/users", func(c *gin.Context) {
		resource, role, err := getResourceAndRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		roleUsers, err := iam.ListResourcesRoleUsers(resource.TenantId, resource.Id, role.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, roleUsers)
	})
	rg.POST("/types/:type/resources/:resource/roles/:role/users", func(c *gin.Context) {
		var roleUser models.ResourceRoleUser
		if c.BindJSON(&roleUser) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		resource, role, err := getResourceAndRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		ru, err := iam.CreateResourceRoleUser(resource.TenantId, resource.Id, role.Id, &roleUser)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ru)
	})
	rg.DELETE("/types/:type/resources/:resource/roles/:role/users/:user", func(c *gin.Context) {
		resource, role, err := getResourceAndRole(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		client, err := GetClientFromCid(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		userName := c.Param("user")
		var clientUser models.ClientUser
		if err := middlewares.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", client.Id, userName); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		var roleUser models.ResourceRoleUser
		if err := middlewares.TenantDB(c).First(&roleUser, "resource_id = ? AND role_id = ? and client_user_id", resource.Id, role.Id, clientUser.Id).Error; err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if err := iam.DeleteResourceRoleUser(resource.TenantId, roleUser.Id); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/types/:type/resources/:resource/actions/:action/users/:user", func(c *gin.Context) {
		resource, err := getResource(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		action, err := getAction(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		client, err := GetClientFromCid(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		userName := c.Param("user")
		var clientUser models.ClientUser
		if err := middlewares.TenantDB(c).First(&clientUser, "client_id = ? AND sub = ?", client.Id, userName); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		result, err := iam.CheckPermission(resource.TenantId, clientUser.Id, resource.Id, action.Id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}

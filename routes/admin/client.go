package admin

import (
	"accounts/data"
	"accounts/middlewares"
	"accounts/models"
	"accounts/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func addAdminClientsRoutes(rg *gin.RouterGroup) {
	rg.GET("/clients", func(c *gin.Context) {
		var clients []models.Client
		if middlewares.TenantDB(c).Find(&clients).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(clients, models.Client2Dto))
	})
	rg.GET("/clients/:clientId", func(c *gin.Context) {
		clientId := c.Param("clientId")
		var client models.Client
		if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, client.Dto())
	})
	rg.POST("/clients", func(c *gin.Context) {
		tenant := middlewares.GetTenant(c)
		var client models.Client
		err := c.BindJSON(&client)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		client.TenantId = tenant.Id
		client.ClientId = uuid.NewString()
		if data.DB.Create(&client).Error != nil {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusOK, client.Dto())
	})
	rg.PUT("/clients/:clientId", func(c *gin.Context) {
		clientId := c.Param("clientId")
		var client models.Client
		if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var cli models.Client
		err := c.BindJSON(&cli)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		client.Name = cli.Name
		if data.DB.Save(&client).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, client.Dto())
	})
	rg.DELETE("/clients/:clientId", func(c *gin.Context) {
		clientId := c.Param("clientId")
		var client models.Client
		if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		if data.DB.Delete(&client).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	rg.GET("/clients/:clientId/redirect-uris", func(c *gin.Context) {
		clientId := c.Param("clientId")
		var client models.Client
		if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var uris []models.RedirectUri
		if middlewares.TenantDB(c).Find(&uris, "client_id = ?", client.Id).Error != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, utils.Filter(uris, models.RedirectUri2Dto))
	})
	rg.POST("/clients/:clientId/redirect-uris", func(c *gin.Context) {
		clientId := c.Param("clientId")
		var client models.Client
		if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		var uri models.RedirectUri
		if c.BindJSON(&uri) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		uri.TenantId = client.TenantId
		uri.ClientId = client.Id
		if middlewares.TenantDB(c).Create(&uri).Error != nil {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusOK, uri.Dto())
	})
	rg.DELETE("/clients/:clientId/redirect-uris/:uriId", func(c *gin.Context) {
		clientId := c.Param("clientId")
		uriId := c.Param("uriId")
		tenant := middlewares.GetTenant(c)
		var uri models.RedirectUri
		if middlewares.TenantDB(c).First(&uri, "tenant_id = ? AND client_id = ? AND uri_id = ?", tenant.Id, clientId, uriId).Error != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusNoContent)
	})
}

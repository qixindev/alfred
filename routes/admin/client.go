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

// ListClients godoc
//
//	@Summary	client
//	@Schemes
//	@Description	get client list
//	@Tags			client
//	@Param			tenant	path		string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients [get]
func ListClients(c *gin.Context) {
	var clients []models.Client
	if middlewares.TenantDB(c).Find(&clients).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(clients, models.Client2Dto))
}

// GetClient godoc
//
//	@Summary	client
//	@Schemes
//	@Description	get client
//	@Tags			client
//	@Param			tenant		path		string	true	"tenant"
//	@Param			clientId	path		integer	true	"clientId"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [get]
func GetClient(c *gin.Context) {
	clientId := c.Param("clientId")
	var client models.Client
	if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, client.Dto())
}

// NewClient godoc
//
//	@Summary	new client
//	@Schemes
//	@Description	new client
//	@Tags			client
//	@Param			tenant	path	string	true	"tenant"
//	@Param			name	body	object	true	"{"name": "main"}"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients [post]
func NewClient(c *gin.Context) {
	tenant := middlewares.GetTenant(c)
	var client models.Client
	if err := c.BindJSON(&client); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	client.TenantId = tenant.Id
	client.CliId = uuid.NewString()
	if data.DB.Create(&client).Error != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, client.Dto())
}

// UpdateClient godoc
//
//	@Summary	update client
//	@Schemes
//	@Description	update client
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"clientId"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [put]
func UpdateClient(c *gin.Context) {
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
}

// DeleteClient godoc
//
//	@Summary	delete client
//	@Schemes
//	@Description	delete client
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"clientId"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [delete]
func DeleteClient(c *gin.Context) {
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
}

// ListClientRedirectUri godoc
//
//	@Summary	get client redirect uris
//	@Schemes
//	@Description	get client redirect uris
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris [get]
func ListClientRedirectUri(c *gin.Context) {
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
}

// NewClientRedirectUri godoc
//
//	@Summary	new client redirect uri
//	@Schemes
//	@Description	new client redirect uri
//	@Tags			client
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/:clientId/redirect-uris [post]
func NewClientRedirectUri(c *gin.Context) {
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
}

// DeleteClientRedirectUri godoc
//
//	@Summary	delete client redirect uris
//	@Schemes
//	@Description	delete client redirect uris
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			uriId		path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris/{uriId} [delete]
func DeleteClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	uriId := c.Param("uriId")
	tenant := middlewares.GetTenant(c)
	var uri models.RedirectUri
	if middlewares.TenantDB(c).First(&uri, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, uriId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := data.DB.Delete(&uri).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetClientSecret godoc
//
//	@Summary	get client secrets
//	@Schemes
//	@Description	get client secrets
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secrets [get]
func GetClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	var client models.Client
	if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var secrets []models.ClientSecret
	if middlewares.TenantDB(c).Find(&secrets, "client_id = ?", client.Id).Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.Filter(secrets, models.ClientSecret2Dto))
}

// NewClientSecret godoc
//
//	@Summary	new client secret
//	@Schemes
//	@Description	new client secret
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secrets [post]
func NewClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	var client models.Client
	if middlewares.TenantDB(c).First(&client, "id = ?", clientId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var secret models.ClientSecret
	if err := c.BindJSON(&secret); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	secret.TenantId = client.TenantId
	secret.ClientId = client.Id
	if middlewares.TenantDB(c).Create(&secret).Error != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, secret.Dto())
}

// DeleteClientSecret godoc
//
//	@Summary	delete client secret
//	@Schemes
//	@Description	delete client secret
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Param			secretId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secret/{secretId} [delete]
func DeleteClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	secretId := c.Param("secretId")
	tenant := middlewares.GetTenant(c)
	var secret models.ClientSecret
	if middlewares.TenantDB(c).First(&secret, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, secretId).Error != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := data.DB.Delete(&secret).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListClientUsers godoc
//
//	@Summary		client user
//	@Schemes
//	@Description	get client user list
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users [get]
func ListClientUsers(c *gin.Context) {
	var clientUser []models.ClientUser
	clientId := c.Param("clientId")
	if err := middlewares.TenantDB(c).Find(&clientUser, "client_id = ?", clientId).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, clientUser)
}

func addAdminClientsRoutes(rg *gin.RouterGroup) {
	rg.GET("/clients", ListClients)
	rg.GET("/clients/:clientId", GetClient)
	rg.POST("/clients", NewClient)
	rg.PUT("/clients/:clientId", UpdateClient)
	rg.DELETE("/clients/:clientId", DeleteClient)
	rg.GET("/clients/:clientId/redirect-uris", ListClientRedirectUri)
	rg.POST("/clients/:clientId/redirect-uris", NewClientRedirectUri)
	rg.DELETE("/clients/:clientId/redirect-uris/:uriId", DeleteClientRedirectUri)
	rg.GET("/clients/:clientId/secrets", GetClientSecret)
	rg.POST("/clients/:clientId/secrets", NewClientSecret)
	rg.DELETE("/clients/:clientId/secret/:secretId", DeleteClientSecret)
	rg.GET("/clients/:clientId/users", ListClientUsers)
}

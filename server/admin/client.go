package admin

import (
	"accounts/global"
	"accounts/models"
	"accounts/server/internal"
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
	if err := internal.TenantDB(c).Find(&clients).Error; err != nil {
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
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, client.Dto())
}

// GetDefaultClient godoc
//
//	@Summary	client
//	@Schemes
//	@Description	get client
//	@Tags			client
//	@Param			tenant		path		string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/default [get]
func GetDefaultClient(c *gin.Context) {
	var client models.Client
	if err := internal.TenantDB(c).First(&client, "name = ?", "default").Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get default client err: " + err.Error())
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
	tenant := internal.GetTenant(c)
	var client models.Client
	if err := c.BindJSON(&client); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	client.TenantId = tenant.Id
	if client.Id == "" {
		client.Id = uuid.NewString()
	}

	if err := global.DB.Create(&client).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("create client err: " + err.Error())
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
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	var cli models.Client
	if err := c.BindJSON(&cli); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	client.Name = cli.Name
	if err := global.DB.Save(&client).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("update clients err: " + err.Error())
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
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	if err := global.DB.Delete(&client).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete client err: " + err.Error())
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
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}

	var uris []models.RedirectUri
	if err := internal.TenantDB(c).Find(&uris, "client_id = ?", client.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get redirect-uris err: " + err.Error())
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
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris [post]
func NewClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	var client models.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	var uri models.RedirectUri
	if err := c.BindJSON(&uri); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	uri.TenantId = client.TenantId
	uri.ClientId = client.Id
	if err := internal.TenantDB(c).Create(&uri).Error; err != nil {
		c.Status(http.StatusConflict)
		return
	}
	c.JSON(http.StatusOK, uri.Dto())
}

// UpdateClientRedirectUri godoc
//
//	@Summary	new client redirect uri
//	@Schemes
//	@Description	new client redirect uri
//	@Tags			client
//	@Param			tenant	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris/{uriId} [post]
func UpdateClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	uriId := c.Param("uriId")
	var newUri models.RedirectUri
	if err := c.BindJSON(&newUri); err != nil {
		internal.ErrReqPara(c, err)
		return
	}

	var uri models.RedirectUri
	if err := internal.TenantDB(c).First(&uri, "client_id = ? AND id = ?", clientId, uriId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}

	uri.RedirectUri = newUri.RedirectUri
	if err := internal.TenantDB(c).Updates(&uri).Error; err != nil {
		c.Status(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newUri.Dto())
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
	tenant := internal.GetTenant(c)
	var uri models.RedirectUri
	if err := internal.TenantDB(c).First(&uri, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, uriId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get redirect uri err: " + err.Error())
		return
	}

	if err := global.DB.Delete(&uri).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete redirect-uri err: " + err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// ListClientSecret godoc
//
//	@Summary	get client secrets
//	@Schemes
//	@Description	get client secrets
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secrets [get]
func ListClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	var client models.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	var secrets []models.ClientSecret
	if err := internal.TenantDB(c).Find(&secrets, "client_id = ?", client.Id).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get clients secret err: " + err.Error())
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
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client err: " + err.Error())
		return
	}
	var secret models.ClientSecret
	if err := c.BindJSON(&secret); err != nil {
		internal.ErrReqPara(c, err)
		return
	}
	secret.TenantId = client.TenantId
	secret.ClientId = client.Id
	if err := internal.TenantDB(c).Create(&secret).Error; err != nil {
		c.Status(http.StatusConflict)
		global.LOG.Error("new client secret err: " + err.Error())
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
	tenant := internal.GetTenant(c)
	var secret models.ClientSecret
	if err := internal.TenantDB(c).First(&secret, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, secretId).Error; err != nil {
		c.Status(http.StatusNotFound)
		global.LOG.Error("get client secret err: " + err.Error())
		return
	}

	if err := global.DB.Delete(&secret).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("delete client secret err: " + err.Error())
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
//	@Param			clientId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users [get]
func ListClientUsers(c *gin.Context) {
	var clientUser []models.ClientUser
	clientId := c.Param("clientId")
	if err := global.DB.Table("client_users cu").Select("cu.sub sub, cu.client_id, u.username user_name").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ?", internal.GetTenant(c).Id, clientId).
		Find(&clientUser).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get client user err: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.Filter(clientUser, models.ClientUserDto))
}

// GetClientUsers godoc
//
//	@Summary		client user
//	@Schemes
//	@Description	get client user list
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"
//	@Param			clientId	path	string	true	"tenant"
//	@Param			subId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users/{subId} [get]
func GetClientUsers(c *gin.Context) {
	var clientUser struct {
		Sub      string `json:"sub"`
		ClientId string `json:"clientId"`
		models.User
	}
	clientId := c.Param("clientId")
	subId := c.Param("subId")
	if err := global.DB.Debug().Table("client_users cu").
		Select("cu.sub sub, cu.client_id, u.username username, u.phone, u.email, u.first_name, u.last_name, u.display_name, u.role").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ? AND cu.sub = ?", internal.GetTenant(c).Id, clientId, subId).
		Find(&clientUser).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		global.LOG.Error("get client user err: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, clientUser)
}

func AddAdminClientsRoutes(rg *gin.RouterGroup) {
	rg.GET("/clients", ListClients)
	rg.GET("/clients/:clientId", GetClient)
	rg.GET("/clients/default", GetDefaultClient)
	rg.POST("/clients", NewClient)
	rg.PUT("/clients/:clientId", UpdateClient)
	rg.DELETE("/clients/:clientId", DeleteClient)
	rg.GET("/clients/:clientId/redirect-uris", ListClientRedirectUri)
	rg.POST("/clients/:clientId/redirect-uris", NewClientRedirectUri)
	rg.PUT("/clients/:clientId/redirect-uris/:uriId", UpdateClientRedirectUri)
	rg.DELETE("/clients/:clientId/redirect-uris/:uriId", DeleteClientRedirectUri)
	rg.GET("/clients/:clientId/secrets", ListClientSecret)
	rg.POST("/clients/:clientId/secrets", NewClientSecret)
	rg.DELETE("/clients/:clientId/secret/:secretId", DeleteClientSecret)
	rg.GET("/clients/:clientId/users", ListClientUsers)
	rg.GET("/clients/:clientId/users/:subId", GetClientUsers)
}

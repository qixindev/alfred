package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
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
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients [get]
func ListClients(c *gin.Context) {
	var clients []model.Client
	if err := internal.TenantDB(c).Find(&clients).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list client err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(clients, model.Client2Dto), 0)
}

// GetClient godoc
//
//	@Summary	client
//	@Schemes
//	@Description	get client
//	@Tags			client
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [get]
func GetClient(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	resp.SuccessWithData(c, client.Dto())
}

// GetDefaultClient godoc
//
//	@Summary	client
//	@Schemes
//	@Description	get client
//	@Tags			client
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/default [get]
func GetDefaultClient(c *gin.Context) {
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "name = ?", "default").Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get default client err")
		return
	}
	resp.SuccessWithData(c, client.Dto())
}

// NewClient godoc
//
//	@Summary	new client
//	@Schemes
//	@Description	new client
//	@Tags			client
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			name	body	object	true	"{"name": "main"}"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients [post]
func NewClient(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var client model.Client
	if err := c.BindJSON(&client); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	client.TenantId = tenant.Id
	if client.Id == "" {
		client.Id = uuid.NewString()
	}

	if err := global.DB.Create(&client).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "failed to create client")
		return
	}

	if err := global.DB.Create(&model.ClientSecret{
		ClientId: client.Id, Name: client.Name, Secret: uuid.NewString(), TenantId: tenant.Id,
	}).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "failed to create client secret")
		return
	}

	resp.SuccessWithData(c, client.Dto())
}

// UpdateClient godoc
//
//	@Summary	update client
//	@Schemes
//	@Description	update client
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [put]
func UpdateClient(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	var cli model.Client
	if err := c.BindJSON(&cli); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	client.Name = cli.Name
	if err := global.DB.Save(&client).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update clients err")
		return
	}
	resp.SuccessWithData(c, client.Dto())
}

// DeleteClient godoc
//
//	@Summary	delete client
//	@Schemes
//	@Description	delete client
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId} [delete]
func DeleteClient(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	tenant := internal.GetTenant(c)
	if err := service.DeleteClient(tenant.Id, clientId); err != nil {
		resp.ErrorSqlDelete(c, err, "delete client err")
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris [get]
func ListClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err", true)
		return
	}

	var uris []model.RedirectUri
	if err := internal.TenantDB(c).Find(&uris, "client_id = ?", client.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list redirect-uris err", true)
		return
	}

	resp.SuccessWithArrayData(c, utils.Filter(uris, model.RedirectUri2Dto), 0)
}

// NewClientRedirectUri godoc
//
//	@Summary	new client redirect uri
//	@Schemes
//	@Description	new client redirect uri
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris [post]
func NewClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	var uri model.RedirectUri
	if err := c.BindJSON(&uri); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	uri.TenantId = client.TenantId
	uri.ClientId = client.Id
	if err := internal.TenantDB(c).Create(&uri).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "create redirect uri err")
		return
	}
	resp.SuccessWithData(c, uri.Dto())
}

// UpdateClientRedirectUri godoc
//
//	@Summary	new client redirect uri
//	@Schemes
//	@Description	new client redirect uri
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris/{uriId} [post]
func UpdateClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	uriId := c.Param("uriId")
	var newUri model.RedirectUri
	if err := c.BindJSON(&newUri); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	var uri model.RedirectUri
	if err := internal.TenantDB(c).First(&uri, "client_id = ? AND id = ?", clientId, uriId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}

	uri.RedirectUri = newUri.RedirectUri
	if err := internal.TenantDB(c).Updates(&uri).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update redirect uri err")
		return
	}

	resp.SuccessWithData(c, newUri.Dto())
}

// DeleteClientRedirectUri godoc
//
//	@Summary	delete client redirect uris
//	@Schemes
//	@Description	delete client redirect uris
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Param			uriId		path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/redirect-uris/{uriId} [delete]
func DeleteClientRedirectUri(c *gin.Context) {
	clientId := c.Param("clientId")
	uriId := c.Param("uriId")
	tenant := internal.GetTenant(c)
	var uri model.RedirectUri
	if err := internal.TenantDB(c).First(&uri, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, uriId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get redirect uri err: ")
		return
	}

	if err := global.DB.Delete(&uri).Error; err != nil {
		resp.ErrorSqlDelete(c, err, "delete redirect uri err")
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
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secrets [get]
func ListClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err: ", true)
		return
	}
	var secrets []model.ClientSecret
	if err := internal.TenantDB(c).Find(&secrets, "client_id = ?", client.Id).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list clients secret err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(secrets, model.ClientSecret2Dto), 0)
}

// NewClientSecret godoc
//
//	@Summary	new client secret
//	@Schemes
//	@Description	new client secret
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secrets [post]
func NewClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	var client model.Client
	if err := internal.TenantDB(c).First(&client, "id = ?", clientId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client err")
		return
	}
	var secret model.ClientSecret
	if err := c.BindJSON(&secret); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	secret.TenantId = client.TenantId
	secret.ClientId = client.Id
	if err := internal.TenantDB(c).Create(&secret).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "create new client secret err")
		return
	}
	resp.SuccessWithData(c, secret.Dto())
}

// DeleteClientSecret godoc
//
//	@Summary	delete client secret
//	@Schemes
//	@Description	delete client secret
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"clientId"	default(default)
//	@Param			secretId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/secret/{secretId} [delete]
func DeleteClientSecret(c *gin.Context) {
	clientId := c.Param("clientId")
	secretId := c.Param("secretId")
	tenant := internal.GetTenant(c)
	var secret model.ClientSecret
	if err := internal.TenantDB(c).First(&secret, "tenant_id = ? AND client_id = ? AND id = ?", tenant.Id, clientId, secretId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get client secret err")
		return
	}

	if err := global.DB.Delete(&secret).Error; err != nil {
		resp.ErrorSqlDelete(c, err, "delete client secret err")
		return
	}

	c.Status(http.StatusNoContent)
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
}

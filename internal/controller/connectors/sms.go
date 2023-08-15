package connectors

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
)

type SmsConnector interface {
	Send(number string, contents []string) error
}

func GetConnector(tenantId, connectorId uint) (SmsConnector, error) {
	var c model.SmsConnector
	if err := global.DB.First(&c, "tenant_id = ? AND id = ?", tenantId, connectorId).Error; err != nil {
		return nil, err
	}
	var connector SmsConnector
	if c.Type == "tcloud" {
		var config model.SmsTcloud
		if err := global.DB.First(&config, "tenant_id = ? AND sms_connector_id = ?", c.TenantId, c.Id).Error; err != nil {
			return nil, err
		}
		sms := SmsTcloud{Config: config}
		connector = &sms
	}
	return connector, nil
}

// GetConnectorList godoc
//
//	@Summary	create connector
//	@Schemes
//	@Description	create connector
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenant}/connectors [get]
func GetConnectorList(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var conn model.SmsConnector
	if err := global.DB.Where("tenant_id = ?", tenant.Id).Find(&conn).Error; err != nil {
		resp.ErrorSqlResponse(c, "")
		return
	}
}

// CreateConnector godoc
//
//	@Summary	create connector
//	@Schemes
//	@Description	create connector
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/tenants/{tenant}/connectors [post]
func CreateConnector(c *gin.Context) {

}

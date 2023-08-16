package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"github.com/gin-gonic/gin"
)

type ControllerRoute struct {
	internal.Api
}

func NewControllerRoute() *ControllerRoute {
	return &ControllerRoute{}
}

// ListConnector godoc
//
//	@Summary	create connector
//	@Schemes
//	@Description	create connector
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/connectors [get]
func (r ControllerRoute) ListConnector(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var conn []model.SmsConnector
	if err := global.DB.Where("tenant_id = ?", tenant.Id).Find(&conn).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "")
		return
	}
	resp.SuccessWithData(c, conn)
}

// GetConnector godoc
//
//	@Summary	create connector
//	@Schemes
//	@Description	create connector
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			connectorId	path	string	true	"connectorId"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/connectors/{connectorId} [get]
func (r ControllerRoute) GetConnector(c *gin.Context) {
	connectorId := c.Param("connectorId")
	tenant := internal.GetTenant(c)
	var conn model.SmsConnector
	if err := global.DB.Where("id = ? AND tenant_id = ?", connectorId, tenant.Id).First(&conn).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get connector err")
		return
	}
	conn.TenantId = tenant.Id
	detail, err := service.GetConnectorDetails(conn)
	if err != nil {
		resp.ErrorUnknown(c, err, "get connector details err")
		return
	}
	resp.SuccessWithData(c, detail)
}

// CreateConnector godoc
//
//	@Summary	create connector
//	@Schemes
//	@Description	create connector
//	@Tags			admin-tenants
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/connectors [post]
func (r ControllerRoute) CreateConnector(c *gin.Context) {

}

// ModifyConnector godoc
//
//	@Summary	modify connector
//	@Schemes
//	@Description	modify connector
//	@Tags			admin-tenants
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			connector	path	string	true	"connectorId"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/connectors/{connectorId} [put]
func (r ControllerRoute) ModifyConnector(c *gin.Context) {

}

func AddConnectorRoute(r *gin.RouterGroup) {
	c := NewControllerRoute()
	r.GET("/connectors", c.ListConnector)
	r.GET("/connectors/:connectorId", c.GetConnector)
	r.POST("/connectors", c.CreateConnector)
	r.PUT("/connectors/:connectorId", c.ModifyConnector)
}

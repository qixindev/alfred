package admin

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/req"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ListSMS
// @Summary	list sms
// @Tags	sms
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Success	200
// @Router	/accounts/admin/{tenant}/sms [get]
func ListSMS(c *gin.Context) {
	var sms []model.SmsConnector
	if err := internal.TenantDB(c).Find(&sms).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list provider err")
		return
	}
	resp.SuccessWithArrayData(c, sms, 0)
}

// GetSMS
// @Summary	get sms
// @Tags	sms
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	smsId	path	integer	true	"sms id"
// @Success	200
// @Router	/accounts/admin/{tenant}/sms/{smsId} [get]
func GetSMS(c *gin.Context) {
	tenant := internal.GetTenant(c)
	smsId := c.Param("smsId")
	var s model.SmsConnector
	if err := internal.TenantDB(c).First(&s, "id = ?", smsId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get sms err")
		return
	}

	res, err := service.GetSmsConfig(tenant.Id, s.Id, s.Type)
	if err != nil {
		resp.ErrorSqlFirst(c, err, "get sms config err")
		return
	}

	resp.SuccessWithData(c, res)
}

// NewSMS
// @Summary	new sms
// @Tags	sms
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	req	body	req.Sms	true	"body"
// @Success	200
// @Router	/accounts/admin/{tenant}/sms [post]
func NewSMS(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var sms req.Sms
	if err := c.BindJSON(&sms); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	sms.Id = 0
	sms.TenantId = tenant.Id
	if err := service.CreateSmsConfig(sms.Type, sms); err != nil {
		resp.ErrorSqlCreate(c, err, "create sms config err")
		return
	}
	resp.Success(c)
}

// UpdateSMS
// @Summary	update sms
// @Tags	sms
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	smsId	path	integer	true	"sms id"
// @Param	req	body	req.Sms	true	"body"
// @Success	200
// @Router	/accounts/admin/{tenant}/sms/{smsId} [put]
func UpdateSMS(c *gin.Context) {
	tenant := internal.GetTenant(c)
	smsId := c.Param("smsId")
	var p req.Sms
	if err := c.BindJSON(&p); err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	p.TenantId = tenant.Id
	p.Id = utils.StrToUint(smsId)
	if err := service.UpdateSmsConfig(p.TenantId, p.Id, p.Type, p); err != nil {
		resp.ErrorSqlUpdate(c, err, "update sms config err")
		return
	}

	resp.Success(c)
}

// DeleteSMS
// @Summary	delete sms
// @Tags	sms
// @Param	tenant	path	string	true	"tenant"	default(default)
// @Param	smsId	path	integer	true	"sms id"
// @Success	200
// @Router	/accounts/admin/{tenant}/sms/{smsId} [delete]
func DeleteSMS(c *gin.Context) {
	smsId := c.Param("smsId")
	var sms model.SmsConnector
	if err := internal.TenantDB(c).First(&sms, "id = ?", smsId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get sms err")
		return
	}

	if err := service.DeleteSmsConfig(sms.TenantId, sms.Id, sms.Type); err != nil {
		resp.ErrorSqlDelete(c, err, "delete sms config err")
		return
	}

	resp.Success(c)
}

func AddAdminSmsRoutes(rg *gin.RouterGroup) {
	rg.GET("/sms", ListSMS)
	rg.GET("/sms/:smsId", GetSMS)
	rg.POST("/sms", NewSMS)
	rg.PUT("/sms/:smsId", UpdateSMS)
	rg.DELETE("/sms/:smsId", DeleteSMS)
}

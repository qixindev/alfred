package admin

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/pkg/global"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// GetLoginPage .
// @Description	get tenant login page
// @Tags		login-page
// @Param		tenant		path	string	true	"租户id"		default(default)
// @Success		200
// @Router		/accounts/admin/{tenant}/clients/{clientId}/page/login [get]
func GetLoginPage(c *gin.Context) {
	tenantName := c.Param("tenant")
	var tenant model.Tenant
	if err := global.DB.Where("name = ?", tenantName).First(&tenant).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	var loginPage map[string]interface{}
	if err := json.Unmarshal([]byte(tenant.LoginPage), &loginPage); err != nil {
		resp.ErrorUnknown(c, err, "unmarshal login page err")
		return
	}
	resp.SuccessWithData(c, loginPage)
}

// UpdateLoginPage .
// @Description	update tenant login page
// @Tags		login-page
// @Param		tenant		path	string	true	"租户名"		default(default)
// @Param		bd			body	object	true	"body"
// @Success		200
// @Router		/accounts/admin/{tenant}/clients/{clientId}/page/login [put]
func UpdateLoginPage(c *gin.Context) {
	var loginPage map[string]interface{}
	if err := internal.New(c).BindJson(&loginPage).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	tenantName := c.Param("tenant")
	var tenant model.Tenant
	if err := global.DB.Where("name = ?", tenantName).First(&tenant).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	pageString, err := json.Marshal(&loginPage)
	if err != nil {
		resp.ErrorUnknown(c, err, "marshal login page err")
		return
	}
	tenant.LoginPage = string(pageString)
	if err = global.DB.Select("login_page").Save(&tenant).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update login page err")
		return
	}

	resp.Success(c)
}

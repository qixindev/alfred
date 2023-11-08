package admin

import (
	"alfred/internal/controller/internal"
	"alfred/internal/endpoint/resp"
	"alfred/internal/model"
	"alfred/internal/service"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetLoginPage .
// @Description	get tenant login page
// @Tags		tenant-meta
// @Param		tenant		path	string	true	"租户"		default(default)
// @Success		200
// @Router		/accounts/admin/{tenant}/page/login [get]
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
// @Tags		tenant-meta
// @Param		tenant		path	string	true	"租户名"		default(default)
// @Param		bd			body	object	true	"body"
// @Success		200
// @Router		/accounts/admin/{tenant}/page/login [put]
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

// GetTenantProto .
// @Description	获取用户隐私协议
// @Tags		tenant-meta
// @Param		tenant		path	string	true	"租户"		default(default)
// @Success		200
// @Router		/accounts/admin/{tenant}/proto [get]
func GetTenantProto(c *gin.Context) {
	tenantName := c.Param("tenant")
	var tenant model.Tenant
	if err := global.DB.Where("name = ?", tenantName).First(&tenant).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	var proto []map[string]interface{}
	if err := json.Unmarshal([]byte(tenant.Proto), &proto); err != nil {
		resp.ErrorUnknown(c, err, "unmarshal tenant proto err")
		return
	}
	resp.SuccessWithData(c, proto)
}

// UpdateTenantProto .
// @Description	update tenant login page
// @Tags		tenant-meta
// @Param		tenant		path	string	true	"租户名"		default(default)
// @Param		bd			body	[]object	true	"body"
// @Success		200
// @Router		/accounts/admin/{tenant}/proto [put]
func UpdateTenantProto(c *gin.Context) {
	var proto []map[string]interface{}
	if err := internal.New(c).BindJson(&proto).Error; err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	tenantName := c.Param("tenant")
	var tenant model.Tenant
	if err := global.DB.Where("name = ?", tenantName).First(&tenant).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get tenant err")
		return
	}
	protoString, err := json.Marshal(&proto)
	if err != nil {
		resp.ErrorUnknown(c, err, "marshal tenant proto err")
		return
	}
	tenant.Proto = string(protoString)
	if err = global.DB.Select("proto").Save(&tenant).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update login page err")
		return
	}

	resp.Success(c)
}

// UploadTenantPicture .
// @Description	上传图片
// @Tags		tenant-meta
// @Param		tenant	path		string	true	"租户名"						default(default)
// @Param		type	path		string	true	"图片类型(background|logo)"	default(logo)
// @Param		file	formData	file	true	"文件"
// @Success		200
// @Router		/accounts/admin/{tenant}/picture/{type}/upload [put]
func UploadTenantPicture(c *gin.Context) {
	tenantName := c.Param("tenant")
	pngType := c.Param("type")
	if pngType != "background" && pngType != "logo" {
		resp.ErrorValidate(c, "pngType does not match")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		resp.ErrorValidate(c, "get file err")
		return
	}

	if file.Size > 2*1024*1024 { // 最大2M
		resp.ErrorValidate(c, "file too large")
		return
	}

	src, err := file.Open()
	if err != nil {
		resp.ErrorUnknown(c, err, "open file error")
		return
	}
	defer utils.DeferErr(src.Close)

	fileName := fmt.Sprintf("alfred-%s-%s.%s", tenantName, pngType, "png")
	url, err := service.UploadFileToAzureBlob(fileName, src)
	if err != nil {
		resp.ErrorUpload(c, err, "upload error")
		return
	}
	resp.SuccessWithData(c, gin.H{"url": url})
}

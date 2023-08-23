package admin

import (
	"accounts/internal/controller/internal"
	"accounts/internal/endpoint/resp"
	"accounts/internal/model"
	"accounts/internal/service"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/url"
	"strings"
)

// ListClientUsers godoc
//
//	@Summary		client user
//	@Schemes
//	@Description	get client user list
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users [get]
func ListClientUsers(c *gin.Context) {
	var clientUser []struct {
		Sub      string `json:"sub"`
		ClientId string `json:"clientId"`
		model.User
	}
	clientId := c.Param("clientId")
	if err := global.DB.Table("client_users cu").
		Select("cu.id, cu.sub sub, cu.client_id, u.username username, u.phone, u.email, u.first_name, u.last_name, u.display_name, u.role").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ?", internal.GetTenant(c).Id, clientId).
		Find(&clientUser).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list client user err", true)
		return
	}
	resp.SuccessWithArrayData(c, clientUser, 0)
}

// GetClientUsers godoc
//
//	@Summary		client user
//	@Schemes
//	@Description	get client user list
//	@Tags			client
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"tenant"
//	@Param			subId		path	string	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users/{subId} [get]
func GetClientUsers(c *gin.Context) {
	var clientUser struct {
		Sub      string `json:"sub"`
		ClientId string `json:"clientId"`
		model.User
	}
	clientId := c.Param("clientId")
	subId := c.Param("subId")
	if err := global.DB.Table("client_users cu").
		Select("cu.id, cu.sub sub, cu.client_id, u.username username, u.phone, u.email, u.first_name, u.last_name, u.display_name, u.role").
		Joins("LEFT JOIN users u ON u.id = cu.user_id").
		Where("cu.tenant_id = ? AND cu.client_id = ? AND cu.sub = ?", internal.GetTenant(c).Id, clientId, subId).
		Find(&clientUser).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "get client user err")
		return
	}

	if clientUser.Username == "" {
		resp.ErrorNotFound(c, "no such client user")
		return
	}
	resp.SuccessWithData(c, clientUser)
}

// UpdateUserPassword godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"tenant"
//	@Param			userId		path	integer	true	"user id"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/password [put]
func UpdateUserPassword(c *gin.Context) {
	var u struct {
		OldPassword         string `json:"oldPassword"`
		NewPassword         string `json:"newPassword"`
		PasswordEncryptType string `json:"passwordEncryptType"`
	}
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if u.NewPassword != u.PasswordEncryptType {
		resp.ErrorRequestWithMsg(c, "PasswordEncrypt failed")
	}

	user, err := service.GetUserBySubId(internal.GetTenant(c).Id, c.Param("clientId"), c.Param("subId"))
	if err != nil {
		resp.ErrorSqlSelect(c, err, "no such user")
	}

	// 检查旧密码
	oldHash, err := utils.HashPassword(u.OldPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}
	if oldHash != user.PasswordHash {
		resp.ErrorRequestWithMsg(c, "invalid old password")
		return
	}

	newHash, err := utils.HashPassword(u.NewPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}
	user.PasswordHash = newHash
	if err = global.DB.Select("password_hash").Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update user password err")
		return
	}
	resp.Success(c)
}

// UpdateUserProfile godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant		path	string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"tenant"
//	@Param			userId		path	integer	true	"user id"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/profile [put]
func UpdateUserProfile(c *gin.Context) {
	var u struct {
		OldPassword         string `json:"oldPassword"`
		NewPassword         string `json:"newPassword"`
		PasswordEncryptType string `json:"passwordEncryptType"`
	}
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if u.NewPassword != u.PasswordEncryptType {
		resp.ErrorRequestWithMsg(c, "PasswordEncrypt failed")
	}

	user, err := service.GetUserBySubId(internal.GetTenant(c).Id, c.Param("clientId"), c.Param("subId"))
	if err != nil {
		resp.ErrorSqlSelect(c, err, "no such user")
	}

	// 检查旧密码
	oldHash, err := utils.HashPassword(u.OldPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}
	if oldHash != user.PasswordHash {
		resp.ErrorRequestWithMsg(c, "invalid old password")
		return
	}

	newHash, err := utils.HashPassword(u.NewPassword)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}
	user.PasswordHash = newHash
	if err = global.DB.Select("password_hash").Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update user password err")
		return
	}
	resp.Success(c)
}

// UpdateAvatar godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant		path		string	true	"tenant"	default(default)
//	@Param			clientId	path	string	true	"tenant"
//	@Param			userId		path		integer	true	"user id"
//	@Param			file		formData	file	true	"file stream"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/clients/{clientId}/users/{subId}/avatar [put]
func UpdateAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		resp.ErrorRequest(c, err)
		return
	}

	fileParts := strings.Split(file.Filename, ".")
	fileType := fileParts[len(fileParts)-1]
	src, err := file.Open()
	if err != nil {
		resp.ErrorUnknown(c, err, "can not open file")
		return
	}
	defer utils.DeferErr(src.Close)

	credential, err := azblob.NewSharedKeyCredential(global.CONFIG.AzureBlob.AccountName, global.CONFIG.AzureBlob.AccountKey)
	if err != nil {
		resp.ErrorUnknown(c, err, "can not upload file")
		return
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, err := url.Parse(global.CONFIG.Urls.AzureBlob)
	if err != nil {
		resp.ErrorUnknown(c, err, "could not parse url")
		return
	}

	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), fileType)
	blobURL := azblob.NewContainerURL(*URL, p).NewBlockBlobURL(fileName)
	if _, err = azblob.UploadStreamToBlockBlob(context.Background(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{}); err != nil {
		resp.ErrorUnknown(c, err, "upload file failed")
		return
	}

	resp.Success(c)
}

func AddClientUserRoute(rg *gin.RouterGroup) {
	rg.GET("/clients/:clientId/users", ListClientUsers)
	rg.GET("/clients/:clientId/users/:subId", GetClientUsers)
	rg.PUT("/clients/:clientId/users/:subId/password", UpdateUserPassword)
	rg.PUT("/clients/:clientId/users/:subId/profile", UpdateUserProfile)
	rg.PUT("/clients/:clientId/users/:subId/avatar", UpdateAvatar)
}

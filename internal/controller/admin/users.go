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

// ListUsers godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user list
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [get]
func ListUsers(c *gin.Context) {
	var users []model.User
	if err := internal.TenantDB(c).Find(&users).Error; err != nil {
		resp.ErrorSqlSelect(c, err, "list tenant users err", true)
		return
	}
	resp.SuccessWithArrayData(c, utils.Filter(users, model.User2AdminDto), 0)
}

// GetUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	get user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [get]
func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// NewUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	new user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			bd		body	object	true	"user body"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users [post]
func NewUser(c *gin.Context) {
	tenant := internal.GetTenant(c)
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	if user.PasswordHash == "" {
		resp.ErrorRequestWithMsg(c, "password should not be null")
		return
	}

	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		resp.ErrorUnknown(c, err, "password hash err")
		return
	}

	user.TenantId = tenant.Id
	user.PasswordHash = hash
	if err = global.DB.Create(&user).Error; err != nil {
		resp.ErrorSqlCreate(c, err, "new tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// UpdateUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"user id"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
	var u model.User
	if err := c.BindJSON(&u); err != nil {
		resp.ErrorRequest(c, err)
		return
	}
	user.Username = u.Username
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.DisplayName = u.DisplayName
	user.Email = u.Email
	user.Phone = u.Phone
	user.Disabled = u.Disabled
	user.Role = u.Role
	if err := global.DB.Save(&user).Error; err != nil {
		resp.ErrorSqlUpdate(c, err, "update tenant user err")
		return
	}
	resp.SuccessWithData(c, user.AdminDto())
}

// UpdateUserPassword godoc
//
//	@Summary	user
//	@Schemes
//	@Description	update user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"user id"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/password [put]
func UpdateUserPassword(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}
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
//	@Param			tenant	path		string	true	"tenant"	default(default)
//	@Param			userId	path		integer	true	"user id"
//	@Param			file	formData	file	true	"file stream"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId}/avatar [put]
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
	_, err = azblob.UploadStreamToBlockBlob(context.Background(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		resp.ErrorUnknown(c, err, "upload file failed")
		return
	}

	resp.Success(c)
}

// DeleteUser godoc
//
//	@Summary	user
//	@Schemes
//	@Description	delete user
//	@Tags			user
//	@Param			tenant	path	string	true	"tenant"	default(default)
//	@Param			userId	path	integer	true	"tenant"
//	@Success		200
//	@Router			/accounts/admin/{tenant}/users/{userId} [delete]
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	var user model.User
	if err := internal.TenantDB(c).First(&user, "id = ?", userId).Error; err != nil {
		resp.ErrorSqlFirst(c, err, "get user err")
		return
	}

	if err := service.DeleteUser(user.Id); err != nil {
		resp.ErrorSqlDelete(c, err, "delete tenant user err")
		return
	}
	resp.Success(c)
}

func AddAdminUsersRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", ListUsers)
	rg.GET("/users/:userId", GetUser)
	rg.POST("/users", NewUser)
	rg.PUT("/users/:userId", UpdateUser)
	rg.PUT("/users/:userId/password", UpdateUserPassword)
	rg.PUT("/users/:userId/avatar", UpdateAvatar)
	rg.DELETE("/users/:userId", DeleteUser)

	rg.GET("/users/:userId/groups", ListUserGroups)
	rg.POST("/users/:userId/groups", NewUserGroup)
	rg.PUT("/users/:userId/groups/:groupId", UpdateUserGroup)
	rg.DELETE("/users/:userId/groups/:groupId", DeleteUserGroup)
}

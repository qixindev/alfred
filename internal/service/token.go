package service

import (
	"alfred/internal/model"
	"alfred/pkg/global"
	"alfred/pkg/utils"
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func getTenant(c *gin.Context) *model.Tenant {
	return c.MustGet("tenant").(*model.Tenant)
}

func GetClientAccessToken(c *gin.Context, client *model.Client) (string, error) {
	tenant := getTenant(c)
	scope := c.Query("scope")
	iss := fmt.Sprintf("%s/%s", utils.GetHostWithScheme(c), tenant.Name)
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["aud"] = []string{client.Id}
	claims["azp"] = client.Id
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["scope"] = scope

	return getToken(tenant.Name, token)
}

func GetAccessToken(c *gin.Context, client *model.Client) (string, error) {
	user := c.MustGet("user").(*model.User)
	tenant := getTenant(c)
	scope := c.Query("scope")
	var clientUser model.ClientUser
	if err := global.DB.First(&clientUser, "tenant_id = ? AND client_id = ? AND user_id = ?", client.TenantId, client.Id, user.Id).Error; err != nil {
		clientUser.TenantId = client.TenantId
		clientUser.ClientId = client.Id
		clientUser.UserId = user.Id
		clientUser.Sub = uuid.NewString()
		if err = global.DB.Create(&clientUser).Error; err != nil {
			return "", err
		}
	}

	iss := fmt.Sprintf("%s/%s", utils.GetHostWithScheme(c), tenant.Name)
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["sub"] = clientUser.Sub
	claims["aud"] = []string{client.Id}
	claims["azp"] = client.Id
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["name"] = user.Name()
	claims["scope"] = scope

	return getToken(tenant.Name, token)
}

// GetResetPasswordToken 生成一次性的重置密码的token
func GetResetPasswordToken(c *gin.Context, phone string) (string, error) {
	tenant := getTenant(c)
	iss := fmt.Sprintf("%s/%s", utils.GetHostWithScheme(c), tenant.Name)
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["sub"] = phone
	claims["exp"] = now.Add(5 * time.Minute).Unix()
	claims["iat"] = now.Unix()

	return getToken(tenant.Name, token)
}

func GetDeviceToken(c *gin.Context, device *model.Device) (string, error) {
	tenant := getTenant(c)
	scope := c.Query("scope")
	iss := fmt.Sprintf("%s/%s", utils.GetHostWithScheme(c), tenant.Name)
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = iss
	claims["aud"] = []string{device.Id}
	claims["azp"] = device.Id
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["scope"] = scope

	return getToken(tenant.Name, token)
}

func getToken(tenant string, token *jwt.Token) (string, error) {
	keys, err := utils.LoadRsaPrivateKeys(tenant)
	if err != nil {
		return "", err
	}

	var kid string
	var key *rsa.PrivateKey
	for kid, key = range keys {
		break
	}

	token.Header["kid"] = kid
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

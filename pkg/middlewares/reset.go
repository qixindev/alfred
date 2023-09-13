package middlewares

import (
	"accounts/internal/endpoint/resp"
	"accounts/pkg/global"
	"accounts/pkg/utils"
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ResetPasswordToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	keys, err := utils.LoadRsaPrivateKeys(c.Param("tenant"))
	if err != nil {
		resp.ErrorUnauthorized(c, err, "load private key err")
		return
	}

	var key *rsa.PrivateKey
	for _, key = range keys {
		claim := jwt.New(jwt.SigningMethodRS256)
		token, err := jwt.ParseWithClaims(tokenString, claim.Claims, func(token *jwt.Token) (interface{}, error) {
			return key.Public(), nil
		})

		subject, err := claim.Claims.GetSubject()
		if err != nil {
			return
		}

		if err == nil && token.Valid {
			c.Set("sub", subject)
			return
		}
		global.LOG.Warn(fmt.Sprintf("%s token valid err: %s", "default", err))
	}

	resp.ErrorUnauthorized(c, nil, "token invalidate")
}

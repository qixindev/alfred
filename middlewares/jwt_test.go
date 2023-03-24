package middlewares

import (
	"accounts/utils"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
	"testing"
	"time"
)

func getToken() (string, *rsa.PrivateKey) {
	now := time.Now()
	token := jwt.New(jwt.SigningMethodRS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}
	claims["iss"] = "a/b"
	claims["sub"] = "jiang_zhao_feng"
	claims["aud"] = "hello"
	claims["azp"] = "world"
	claims["exp"] = now.Add(24 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["scope"] = ""
	pemString, err := os.ReadFile("../config/jwks/default/4e44009c-b74d-406f-90dc-c84cb1ac67d1.key")
	if err != nil {
		return "", nil
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemString)
	if err != nil {
		return "", nil
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", nil
	}

	return tokenString, key
}

func TestPhaseToken(t *testing.T) {
	tokenString, key := getToken()
	if tokenString == "" || key == nil {
		return
	}

	claim := jwt.New(jwt.SigningMethodRS256)
	token, err := jwt.ParseWithClaims(tokenString, claim.Claims, func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	})

	fmt.Println(token.Valid, err, utils.StructToString(claim))
}

func TestValidateToken(t *testing.T) {
	tokenString, key := getToken()
	if tokenString == "" || key == nil {
		return
	}

	parts := strings.Split(tokenString, ".")
	if err := jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key.Public()); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("signature verified")
}

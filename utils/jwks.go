package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"strings"
)

func LoadRsaPrivateKeys(tenant string) (map[string]*rsa.PrivateKey, error) {
	path := "config/jwks/" + tenant
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]*rsa.PrivateKey)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".key") == false {
			continue
		}
		name := path + "/" + file.Name()
		pemString, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM(pemString)
		if err != nil {
			return nil, err
		}
		kid := strings.Split(file.Name(), ".")[0]
		keys[kid] = key
	}
	return keys, nil
}

func LoadKeys(tenant string) (*jose.JSONWebKeySet, error) {
	path := "config/jwks/" + tenant
	_, err := os.ReadDir(path)
	if err != nil {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return nil, err
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
		writeFile := fmt.Sprintf("%s/%s.key", path, uuid.New())
		payload := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		err = os.WriteFile(writeFile, payload, 0400)
		if err != nil {
			return nil, err
		}
	}

	var jwks jose.JSONWebKeySet
	privateKeys, err := LoadRsaPrivateKeys(tenant)
	if err != nil {
		return nil, err
	}
	for kid, key := range privateKeys {
		pub := key.Public()
		jwk := jose.JSONWebKey{
			Key:                       pub,
			KeyID:                     kid,
			Algorithm:                 "RS256",
			Use:                       "sig",
			Certificates:              []*x509.Certificate{},
			CertificateThumbprintSHA1: []uint8{},
		}
		jwks.Keys = append(jwks.Keys, jwk)
	}

	return &jwks, nil
}

package initial

import (
	"alfred/backend/pkg/config/env"
	"alfred/backend/pkg/utils"
	"errors"
	"fmt"
	"os"
)

const (
	SecretKey     = "accounts-secret"
	DefaultSecret = "secret"
)

func getConfigmapSecret() ([]byte, error) {
	sClient, err := utils.GetK8sClient()
	if err != nil {
		return nil, err
	}

	cm, err := utils.GetConfigMap(sClient, env.DefaultCmJWKS)
	if err != nil || cm.Data == nil {
		return nil, err
	}

	if cm.Data == nil || cm.Data[SecretKey] == "" {
		return nil, errors.New("secret is nil")
	}

	return []byte(cm.Data[SecretKey]), nil
}

func getSecretFile() ([]byte, error) {
	str, err := os.ReadFile("./data/" + SecretKey)
	if err != nil || str == nil {
		return nil, err
	}

	return str, nil
}

func GetSessionSecret() (res []byte) {
	var err error
	if env.GetDeployType() == "k8s" {
		res, err = getConfigmapSecret()
	} else {
		res, err = getSecretFile()
	}

	if err != nil || res == nil {
		fmt.Println("get secret err use default: ", err)
		return []byte(DefaultSecret)
	}

	return res
}

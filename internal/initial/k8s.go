package initial

import (
	"accounts/pkg/config"
	"accounts/pkg/config/env"
	"accounts/utils"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CMKeyZap        = "zap"
	CMKeySys        = "system"
	CMKeyPgsql      = "pgsql"
	CMKeyTencentCos = "tencent-cos"
	CMKeyAliyunOss  = "aliyun-oss"
	CMKeyAzureBlob  = "azure-blob"
	CMKeyRabbitMq   = "rabbit-mq"
)

func GetK8sConfig() (*config.Config, error) {
	conf := config.Config{}
	cm, err := GetConfigMap()
	if err != nil {
		return nil, err
	}

	if err = utils.GetAnyString(&(conf.Zap), cm[CMKeyZap]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.System), cm[CMKeySys]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.Pgsql), cm[CMKeyPgsql]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.TencentCOS), cm[CMKeyTencentCos]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.AliyunOSS), cm[CMKeyAliyunOss]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.AzureBlob), cm[CMKeyAzureBlob]); err != nil {
		return nil, err
	}
	if err = utils.GetAnyString(&(conf.RabbitMq), cm[CMKeyRabbitMq]); err != nil {
		return nil, err
	}

	return &conf, nil
}

func GetConfigMap() (map[string]string, error) {
	sClient, err := utils.GetK8sClient()
	if err != nil {
		return nil, err
	}

	configMaps, err := sClient.CoreV1().ConfigMaps(env.GetNameSpace()).
		Get(context.TODO(), env.GetServiceConfigMapName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMaps.Data, nil
}

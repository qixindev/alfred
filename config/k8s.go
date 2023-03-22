package config

import (
	"accounts/config/env"
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
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

func GetAnyString(a any, s string) error {
	if err := json.Unmarshal([]byte(s), a); err != nil {
		return err
	}

	return nil
}

func GetK8sClient() (*kubernetes.Clientset, error) {
	var kubeConfigStr string
	if home := homedir.HomeDir(); home != "" {
		kubeConfigStr = filepath.Join(home, ".kube", "config")
	} else {
		kubeConfigStr = env.K8sConfigPath
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigStr)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func GetK8sConfig() (*Config, error) {
	conf := Config{}
	cm, err := GetConfigMap()
	if err != nil {
		return nil, err
	}

	if err = GetAnyString(&(conf.Zap), cm[CMKeyZap]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.System), cm[CMKeySys]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.Pgsql), cm[CMKeyPgsql]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.TencentCOS), cm[CMKeyTencentCos]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.AliyunOSS), cm[CMKeyAliyunOss]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.AzureBlob), cm[CMKeyAzureBlob]); err != nil {
		return nil, err
	}
	if err = GetAnyString(&(conf.RabbitMq), cm[CMKeyRabbitMq]); err != nil {
		return nil, err
	}

	return &conf, nil
}

func GetConfigMap() (map[string]string, error) {
	sClient, err := GetK8sClient()
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

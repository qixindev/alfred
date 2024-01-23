package initial

import (
	"alfred/backend/pkg/config"
	"alfred/backend/pkg/config/env"
	"alfred/backend/pkg/global"
	"alfred/backend/pkg/utils"
	"context"
	"fmt"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitConfig() (err error) {
	if env.GetDeployType() == "k8s" {
		global.CONFIG, err = GetK8sConfig()
	} else {
		global.CONFIG, err = Viper(env.GetConfigPath(), "yaml")
	}

	return err
}

func Viper(path string, t string) (*config.Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(t)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Fatal error config file: " + err.Error())
		return nil, err
	}

	conf := config.Config{}
	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println("unmarshal conf err: " + err.Error())
		return nil, err
	}

	return &conf, nil
}

const (
	CMKeyZap        = "zap"
	CMKeySys        = "system"
	CMKeyPgsql      = "pgsql"
	CMKeyTencentCos = "tencent-cos"
	CMKeyAliyunOss  = "aliyun-oss"
	CMKeyAzureBlob  = "azure-blob"
	CMKeyRabbitMq   = "rabbit-mq"
	CMUrls          = "urls"
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
	if err = utils.GetAnyString(&(conf.Urls), cm[CMUrls]); err != nil {
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

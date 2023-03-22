package config

import (
	"accounts/config/env"
	"accounts/global"
	"accounts/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Zap        *Zap        `mapstructure:"zap" json:"zap" yaml:"zap"`
	Pgsql      *Pgsql      `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	System     *System     `mapstructure:"system" json:"system" yaml:"system"`
	AliyunOSS  *AliyunOSS  `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	TencentCOS *TencentCOS `mapstructure:"tencent-cos" json:"tencent-cos" yaml:"tencent-cos"`
	AzureBlob  *AzureBlob  `mapstructure:"azure-blob" json:"azure-blob" yaml:"azure-blob"`
	RabbitMq   *RabbitMq   `mapstructure:"rabbit-mq" json:"rabbit-mq" yaml:"rabbit-mq"`
}

func Viper(path string, t string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(t)
	if err := v.ReadInConfig(); err != nil {
		global.LOG.Error("Fatal error config file: " + err.Error())
		return nil, err
	}

	conf := Config{}
	if err := v.Unmarshal(&conf); err != nil {
		global.LOG.Error("unmarshal conf err: " + err.Error())
		return nil, err
	}

	return &conf, nil
}

func InitConfig() (err error) {
	if env.GetDeployType() == "k8s" {
		global.CONFIG, err = utils.GetK8sConfig()
	} else {
		global.CONFIG, err = Viper(env.GetConfigPath(), "yaml")
	}

	return err
}

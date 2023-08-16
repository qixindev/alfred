package initial

import (
	"accounts/pkg/config"
	"accounts/pkg/config/env"
	"accounts/pkg/global"
	"fmt"
	"github.com/spf13/viper"
)

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

func InitConfig() (err error) {
	if env.GetDeployType() == "k8s" {
		global.CONFIG, err = GetK8sConfig()
	} else {
		global.CONFIG, err = Viper(env.GetConfigPath(), "yaml")
	}

	return err
}

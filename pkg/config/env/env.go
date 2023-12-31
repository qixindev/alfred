package env

import "os"

const (
	K8sConfigPath           = "/root/.kube/config"
	DefaultNameSpace        = "default"
	DefaultServiceConfigMap = "service-config"
	DefaultConfigPath       = "config.dev.yaml"
	DefaultDeployType       = "k8s"
	DefaultCmJWKS           = "accounts.jwks"
	DefaultReleaseType      = "release"
)

func getEnv(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		value = defaultValue
	}
	return value
}

func GetDeployType() string {
	return getEnv("DEPLOY_TYPE", DefaultDeployType)
}

func GetConfigPath() string {
	return getEnv("CONFIG_PATH", DefaultConfigPath)
}

func GetNameSpace() string {
	return getEnv("NANE_SPACE", DefaultNameSpace)
}

func GetServiceConfigMapName() string {
	return getEnv("SERVICE_CONFIG_MAP", DefaultServiceConfigMap)
}

func GetReleaseType() string {
	return getEnv("RELEASE_TYPE", DefaultReleaseType)
}

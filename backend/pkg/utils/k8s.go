package utils

import (
	"alfred/backend/pkg/config/env"
	"context"
	"encoding/json"
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
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

func CreateConfigMap(sClient *kubernetes.Clientset, metaName string, data map[string]string) (*v1.ConfigMap, error) {
	configMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: metaName,
		},
		Data: data,
	}

	cm, err := sClient.CoreV1().ConfigMaps(env.GetNameSpace()).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func UpdateConfigMap(sClient *kubernetes.Clientset, cm *v1.ConfigMap) error {
	if sClient == nil || cm == nil {
		return errors.New("client or configmap is nil")
	}

	if _, err := sClient.CoreV1().ConfigMaps(env.GetNameSpace()).Update(context.TODO(), cm, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func GetConfigMap(sClient *kubernetes.Clientset, metaName string) (*v1.ConfigMap, error) {
	configMaps, err := sClient.CoreV1().ConfigMaps(env.GetNameSpace()).Get(context.TODO(), metaName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMaps, nil
}

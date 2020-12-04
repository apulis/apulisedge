// Copyright 2020 Apulis Technology Inc. All rights reserved.

package utils

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var logger = loggers.LogInstance()

// replace this with the K8s Master IP
var KubeMaster string
var Kubeconfig string
var KubeQPS = float32(5.000000)
var KubeBurst = 10
var KubeContentType = "application/vnd.kubernetes.protobuf"

func InitKubeClient(kubeMaster string, kubeConfig string) {
	KubeMaster = kubeMaster
	Kubeconfig = kubeConfig
	logger.Infof("kubeMaster = %s, kubeConfigPath = %s", KubeMaster, Kubeconfig)
}

// KubeConfig from flags
func KubeConfig() (conf *rest.Config, err error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(KubeMaster, Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeConfig.QPS = KubeQPS
	kubeConfig.Burst = KubeBurst
	kubeConfig.ContentType = KubeContentType
	return kubeConfig, err
}

func GetNodeClient() (corev1.NodeInterface, error) {
	kubeConfig, err := KubeConfig()
	if err != nil {
		logger.Error("Failed to create KubeConfig , error : %v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Error("Failed to create clientset , error : %v", err)
		return nil, err
	}

	return clientSet.CoreV1().Nodes(), nil
}

func GetDeploymentClient(namespace string) (appsv1.DeploymentInterface, error) {
	kubeConfig, err := KubeConfig()
	if err != nil {
		logger.Error("Failed to create KubeConfig , error : %v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Error("Failed to create clientset , error : %v", err)
		return nil, err
	}

	return clientSet.AppsV1().Deployments(namespace), nil
}

func ListNodes() (result *v1.NodeList, err error) {
	nodeClient, err := GetNodeClient()
	if err != nil {
		return nil, err
	}

	result, err = nodeClient.List(context.Background(), metav1.ListOptions{})
	return result, err
}

func DescribeNode(name string) (result *v1.Node, err error) {
	nodeClient, err := GetNodeClient()
	if err != nil {
		return nil, err
	}

	result, err = nodeClient.Get(context.Background(), name, metav1.GetOptions{})
	return result, err
}

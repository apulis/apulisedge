// Copyright 2020 Apulis Technology Inc. All rights reserved.

package utils

import (
	"context"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var logger = loggers.LogInstance()

// replace this with the K8s Master IP
var KubeMaster = "https://cls-f4x353m8.ccs.tencent-cloud.com"
var Kubeconfig = "/.kube/config"
var KubeQPS = float32(5.000000)
var KubeBurst = 10
var KubeContentType = "application/vnd.kubernetes.protobuf"

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

func GetNodeClient() typedv1.NodeInterface {
	kubeConfig, err := KubeConfig()
	if err != nil {
		logger.Error("Failed to create KubeConfig , error : %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Error("Failed to create clientset , error : %v", err)
	}

	return clientSet.CoreV1().Nodes()
}

func ListNodes() (result *v1.NodeList, err error) {
	nodeClient := GetNodeClient()
	result, err = nodeClient.List(context.Background(), metav1.ListOptions{})
	return result, err
}

func DescribeNode(name string) (result *v1.Node, err error) {
	nodeClient := GetNodeClient()
	result, err = nodeClient.Get(context.Background(), name, metav1.GetOptions{})
	return result, err
}

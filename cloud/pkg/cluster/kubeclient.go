// Copyright 2020 Apulis Technology Inc. All rights reserved.

package cluster

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// init kubeclient
func (c *Cluster) InitKube(kubeMaster string, kubeConfig string) {
	c.KubeMaster = kubeMaster
	c.KubeConfFile = kubeConfig
	c.KubeQPS = float32(5.000000)
	c.KubeBurst = 10
	c.KubeContentType = "application/vnd.kubernetes.protobuf"

	logger.Infof("kubeMaster = %s, kubeConfigPath = %s", c.KubeMaster, c.KubeConfFile)
}

// KubeConfig from flags
func (c *Cluster) KubeConfig() (conf *rest.Config, err error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(c.KubeMaster, c.KubeConfFile)
	if err != nil {
		return nil, err
	}
	kubeConfig.QPS = c.KubeQPS
	kubeConfig.Burst = c.KubeBurst
	kubeConfig.ContentType = c.KubeContentType
	return kubeConfig, err
}

func (c *Cluster) GetNodeClient() (corev1.NodeInterface, error) {
	kubeConfig, err := c.KubeConfig()
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

func (c *Cluster) GetDeploymentClient(namespace string) (appsv1.DeploymentInterface, error) {
	kubeConfig, err := c.KubeConfig()
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

func (c *Cluster) GetPodClient(namespace string) (corev1.PodInterface, error) {
	kubeConfig, err := c.KubeConfig()
	if err != nil {
		logger.Error("Failed to create KubeConfig , error : %v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Error("Failed to create clientset , error : %v", err)
		return nil, err
	}

	return clientSet.CoreV1().Pods(namespace), nil
}

func (c *Cluster) ListNodes() (result *v1.NodeList, err error) {
	nodeClient, err := c.GetNodeClient()
	if err != nil {
		return nil, err
	}

	result, err = nodeClient.List(context.Background(), metav1.ListOptions{})
	return result, err
}

func (c *Cluster) DescribeNode(name string) (result *v1.Node, err error) {
	nodeClient, err := c.GetNodeClient()
	if err != nil {
		return nil, err
	}

	result, err = nodeClient.Get(context.Background(), name, metav1.GetOptions{})
	return result, err
}

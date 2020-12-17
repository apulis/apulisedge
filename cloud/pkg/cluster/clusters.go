// Copyright 2020 Apulis Technology Inc. All rights reserved.

package cluster

import (
	"errors"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
)

var (
	ErrFindCluster = errors.New("cluster not exist")
)

var logger = loggers.LogInstance()

var clusters []Cluster

type Cluster struct {
	Id              int64
	Domain          string
	KubeMaster      string
	KubeConfFile    string
	KubeQPS         float32
	KubeBurst       int
	KubeContentType string
	HarborAddress   string
	HarborProject   string
	HarborUser      string
	HarborPasswd    string
	DownloadAddress string
}

func GetCluster(clusterId int64) (*Cluster, error) {
	for _, v := range clusters {
		if v.Id == clusterId {
			return &v, nil
		}
	}
	return nil, ErrFindCluster
}

func InitClusters(config *configs.EdgeCloudConfig) {
	for _, c := range config.Clusters {
		if c.Id < 0 || c.Domain == "" || c.KubeMaster == "" || c.KubeConfFile == "" ||
			c.HarborAddress == "" || c.HarborProject == "" || c.HarborUser == "" ||
			c.HarborPasswd == "" || c.DownloadAddress == "" {
			logger.Errorf("Invalid cluster config = %v", c)
			continue
		}
		clu := Cluster{}
		clu.Domain = c.Domain
		clu.InitKube(c.KubeMaster, c.KubeConfFile)
		clu.InitDockerCli(c.HarborAddress, c.HarborProject, c.HarborUser, c.HarborPasswd)
		clu.DownloadAddress = c.DownloadAddress
		clusters = append(clusters, clu)
		logger.Infof("Add cluster = %v", clu)
	}
}

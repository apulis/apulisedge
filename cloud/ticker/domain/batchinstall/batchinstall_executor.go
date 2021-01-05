package batchinstall

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	nodeservice "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/service"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"golang.org/x/crypto/ssh"
)

var logger = loggers.LogInstance()

func CreateBatchInstallTicker(ctx context.Context, config *configs.EdgeCloudConfig) {

	duration := time.Duration(1) * time.Second
	// duration := time.Duration(config.Portal.NodeCheckerInterval) * time.Second
	installerTicker := time.NewTimer(duration)
	defer installerTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("CreateNodeTickerLoop was terminated")
			return
		case <-installerTicker.C:
			InstallerTicker(config)
			installerTicker.Reset(duration)
		}
	}
}

func InstallerTicker(config *configs.EdgeCloudConfig) {
	var nodeInfos []nodeentity.NodeBasicInfo

	database.Db.Where("Status = ?", constants.StatusInstalling).
		Limit(10).
		Find(&nodeInfos)
	var installingNodeCount int64
	database.Db.Model(&nodeentity.NodeBasicInfo{}).Where("Status = ?", constants.StatusInstalling).Count(&installingNodeCount)
	logger.Infoln("Remain installing node :", installingNodeCount)
	for _, nodeinfo := range nodeInfos {
		logger.Infoln("Ready to install.")
		var batchNode nodeentity.NodeOfBatchInfo
		database.Db.Where("NodeID = ?", nodeinfo.ID).Find(&batchNode)
		installScript, err := nodeservice.CreateInstallScripts(nodeinfo, nodeinfo.Arch)
		if err != nil {
			logger.Debugln("CreateInstallScripts failed.")
			logger.Errorln(err)
			continue
		}
		logger.Infoln("Install node:", batchNode)
		logger.Infoln(installScript)
		sshHost := batchNode.Address
		sshPort, err := strconv.Atoi(batchNode.Port)
		if err != nil {
			logger.Debugln(err)
			continue
		}
		sshUser := batchNode.Sudoer
		sshPassword := batchNode.Password
		sshConfig := &ssh.ClientConfig{
			Config:          ssh.Config{},
			User:            sshUser,
			Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			ClientVersion:   "",
			Timeout:         time.Minute,
		}
		addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
		sshClient, err := ssh.Dial("tcp", addr, sshConfig)
		if err != nil {
			logger.Debugln("Fail to create ssh client.")
			logger.Debugln(err)
			continue
		}
		sshSession, err := sshClient.NewSession()
		if err != nil {
			logger.Debugln("Fail to create ssh session.")
			logger.Debugln(err)
			continue
		}
		combo, err := sshSession.CombinedOutput(installScript)
		if err != nil {
			logger.Debugln("Fail to execute script.")
			logger.Debugln(err)
			continue
		}
		logger.Infoln(string(combo))

		sshClient.Close()
		sshSession.Close()

		database.Db.Model(&nodeinfo).Update("Status", constants.StatusOffline)
	}
}

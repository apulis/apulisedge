package nodeservice

import (
	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

// GetInstallScripts generate install script
func GetInstallScripts(userInfo proto.ApulisHeader, req *nodemodule.GetInstallScriptReq) (string, error) {
	var err error
	var script string
	var nodeInfo nodeentity.NodeBasicInfo

	targetArch := req.Arch

	// get node
	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?",
			userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
		First(&nodeInfo)
	if res.Error != nil {
		return "", res.Error
	}

	// check type
	archExist := true
	if !cluster.IsArchValid(req.Arch) {
		archExist = false
	}

	if !archExist {
		return "", cluster.ErrArchTypeNotExist
	}

	script, err = CreateInstallScripts(nodeInfo, targetArch)
	if err != nil {
		return "", err
	}

	return script, err
}

func CreateInstallScripts(nodeInfo nodeentity.NodeBasicInfo, targetArch string) (string, error) {
	var script string

	var downloadTarget = "/tmp/apulisedge/"
	var packageName = "apulisedge_" + targetArch
	var fileName = packageName + ".tar.gz"
	var pubKeyFileName = packageName + ".key"
	var signFileName = packageName + ".sig"
	var kubeedgeImageName = "apulisedge/apulis/kubeedge-edge:1.0-" + targetArch

	clu, err := cluster.GetCluster(nodeInfo.ClusterId)
	if err != nil {
		logger.Infof("GetInstallScripts, can`t find cluster %d", nodeInfo.ClusterId)
		return "", err
	}

	// clean env
	script = "rm -rf " + downloadTarget + " && "
	// make neccessary dir
	script = script + " mkdir -p " + downloadTarget + " && "
	script = script + " mkdir -p /opt/apulisedge && "
	// download package and signature
	script = script + "wget " + clu.DownloadAddress + "/" + fileName + " -P " + downloadTarget + " && "
	script = script + "wget " + clu.DownloadAddress + "/" + signFileName + " -P " + downloadTarget + " && "
	script = script + "wget " + clu.DownloadAddress + "/" + pubKeyFileName + " -P " + downloadTarget + " && "
	// verify file
	script = script + " openssl dgst -verify " + downloadTarget + "/" + pubKeyFileName + " -sha256 -signature " + downloadTarget + "/" + signFileName + " " + downloadTarget + "/" + fileName + " && "
	// decompress package
	script = script + "tar -zxvf " + downloadTarget + "/" + fileName + " -C " + downloadTarget + " && "
	// move install script
	script = script + "cp " + downloadTarget + "/package/scripts/* /opt/apulisedge/" + " && "
	// run install script
	script = script + "/opt/apulisedge/install_edge.sh -d " + clu.Domain +
		" -l " + clu.HarborAddress + "/" + kubeedgeImageName +
		" -h " + nodeInfo.UniqueName

	return script, nil
}

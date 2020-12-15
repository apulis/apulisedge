package nodeservice

import (
	"strconv"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
)

var DownloadAddress string
var DownloadPort int
var CloudServer string
var ImageServer string

// InitInstallScriptConfig init config
func InitInstallScriptConfig(config *configs.EdgeCloudConfig) {
	DownloadAddress = config.ScriptConfig.DownloadAddress
	DownloadPort = config.ScriptConfig.DownloadPort
	CloudServer = config.ScriptConfig.CloudServer
	ImageServer = config.ScriptConfig.ImageServer
}

// GetInstallScripts generate install script
func GetInstallScripts(req *nodemodule.GetInstallScriptReq) (string, error) {
	var err error
	var script string
	var targetArch = req.Arch
	var downloadTarget = "/tmp/apulisedge/"
	var packageName = "apulisedge_" + targetArch
	var fileName = packageName + ".tar.gz"
	var pubKeyFileName = packageName + ".key"
	var signFileName = packageName + ".sig"

	// clean env
	script = "rm -rf " + downloadTarget + " && "
	// make neccessary dir
	script = script + " mkdir -p " + downloadTarget + " && "
	script = script + " mkdir -p /opt/apulisedge && "
	// download package and signature
	script = script + "wget " + DownloadAddress + ":" + strconv.Itoa(DownloadPort) + "/" + fileName + " -P " + downloadTarget + " && "
	script = script + "wget " + DownloadAddress + ":" + strconv.Itoa(DownloadPort) + "/" + signFileName + " -P " + downloadTarget + " && "
	script = script + "wget " + DownloadAddress + ":" + strconv.Itoa(DownloadPort) + "/" + pubKeyFileName + " -P " + downloadTarget + " && "
	// verify file
	script = script + " openssl dgst -verify " + downloadTarget + "/" + pubKeyFileName + " -sha256 -signature " + downloadTarget + "/" + signFileName + " " + downloadTarget + "/" + fileName + " && "
	// decompress package
	script = script + "tar -zxvf " + downloadTarget + "/" + fileName + " -C " + downloadTarget + " && "
	// move install script
	script = script + "cp " + downloadTarget + "/package/scripts/* /opt/apulisedge/" + " && "
	// run install script
	script = script + "/opt/apulisedge/install_edge.sh -d " + CloudServer + " -l " + ImageServer + "/apulisedge/apulis/kubeedge-edge:1.0"

	return script, err
}

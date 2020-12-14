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
	var fileName = "kubeedgeRuntime-" + targetArch + ".tar.gz"

	// make download dir
	script = " mkdir -p " + downloadTarget + " && "
	// download package
	script = script + "wget " + DownloadAddress + ":" + strconv.Itoa(DownloadPort) + "/" + fileName + " -P " + downloadTarget + " && "
	// decompress package
	script = script + "tar " + downloadTarget + "/" + fileName + " && "
	// move install script
	script = script + "tar " + downloadTarget + "/" + fileName + " && "
	// run install script
	script = script + "/opt/apulisedge/install_edge.sh -d " + CloudServer + " -l " + ImageServer + "/apulisedge/apulis/kubeedge-edge:1.0"

	return script, err
}

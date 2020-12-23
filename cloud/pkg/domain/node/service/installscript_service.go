package nodeservice

import (
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
)

// GetInstallScripts generate install script
func GetInstallScripts(req *nodemodule.GetInstallScriptReq, domain string, imgServer string, dwAddress string) (string, error) {
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
	script = script + "wget " + dwAddress + "/" + fileName + " -P " + downloadTarget + " && "
	script = script + "wget " + dwAddress + "/" + signFileName + " -P " + downloadTarget + " && "
	script = script + "wget " + dwAddress + "/" + pubKeyFileName + " -P " + downloadTarget + " && "
	// verify file
	script = script + " openssl dgst -verify " + downloadTarget + "/" + pubKeyFileName + " -sha256 -signature " + downloadTarget + "/" + signFileName + " " + downloadTarget + "/" + fileName + " && "
	// decompress package
	script = script + "tar -zxvf " + downloadTarget + "/" + fileName + " -C " + downloadTarget + " && "
	// move install script
	script = script + "cp " + downloadTarget + "/package/scripts/* /opt/apulisedge/" + " && "
	// run install script
	script = script + "/opt/apulisedge/install_edge.sh -d " + domain + " -l " + imgServer + "/apulisedge/apulis/kubeedge-edge:1.0"

	return script, err
}

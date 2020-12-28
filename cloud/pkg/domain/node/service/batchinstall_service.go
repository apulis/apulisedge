package nodeservice

import (
	"encoding/csv"
	"io/ioutil"
	"strings"
	"time"

	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

func AnalyzeCSV(userInfo proto.ApulisHeader) error {
	filein, err := ioutil.ReadFile(nodemodule.CSVSavePath)
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	r := csv.NewReader(strings.NewReader(string(filein)))
	records, err := r.ReadAll()
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	for _, line := range records {
		logger.Infoln(line)
	}

	return nil
}

func ListBatchList(userInfo proto.ApulisHeader) (*[]nodeentity.NodeOfBatchInfo, error) {
	var batchInfos []nodeentity.NodeOfBatchInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Find(&batchInfos)

	if res.Error != nil {
		return &batchInfos, res.Error
	}

	return &batchInfos, nil
}

func UpdateBatch(userInfo proto.ApulisHeader) error {
	var batchInfos []nodeentity.NodeOfBatchInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId).
		Find(&batchInfos)
	if res.Error != nil {
		return res.Error
	}

	var err error
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		return err
	}
	for _, batchInfo := range batchInfos {
		uniqName := clu.GetUniqueName(cluster.ResourceNode)
		node := &nodeentity.NodeBasicInfo{
			NodeName:         batchInfo.NodeName,
			ClusterId:        batchInfo.ClusterId,
			GroupId:          batchInfo.GroupId,
			UserId:           batchInfo.UserId,
			NodeType:         batchInfo.NodeType,
			Arch:             batchInfo.Arch,
			UniqueName:       uniqName,
			CpuCore:          0,
			Memory:           0,
			Status:           constants.StatusInstalling,
			Roles:            "",
			ContainerRuntime: "",
			OsImage:          "",
			InterIp:          "",
			OuterIp:          "",
			CreateAt:         time.Now(),
			UpdateAt:         time.Now(),
		}
		err = nodeentity.CreateNode(node)
		if err != nil {
			return err
		}
	}

	return nil
}
func DeleteNodeOfBatch(userInfo proto.ApulisHeader, req *nodemodule.DeleteNodeOfBatchReq) error {
	var nodeOfBatchInfo nodeentity.NodeOfBatchInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and NodeName = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.Name).
		First(&nodeOfBatchInfo)
	if res.Error != nil {
		return res.Error
	}

	return nodeentity.DeleteNodeOfBatch(&nodeOfBatchInfo)
}

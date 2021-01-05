package nodeservice

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/apulis/ApulisEdge/cloud/pkg/cluster"
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	constants "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

func init() {
	_, err := os.Stat(nodemodule.CSVSavePath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(nodemodule.CSVSavePath, 0755)
		} else {
			logger.Panicln("Fatal create file directory error: %s", err)
		}
	}
}

func CreateBatchTask(userInfo proto.ApulisHeader) (int, error) {
	batchTask := &nodeentity.BatchTaskRecord{
		ClusterId:      userInfo.ClusterId,
		GroupId:        userInfo.GroupId,
		UserId:         userInfo.UserId,
		Status:         "parsing",
		ParsingFileNum: 0,
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
	}
	err := nodeentity.CreateBatchTask(batchTask)
	if err != nil {
		return 0, err
	}
	return int(batchTask.ID), nil
}

func AnalyzeCSV(batchID int, filePath string) error {
	nodeentity.BeginBatchTask(batchID)
	go csvParseTask(batchID, filePath)

	return nil
}

func csvParseTask(batchID int, filePath string) error {
	taskInfo, err := nodeentity.GetBatchTaskByID(batchID)
	if err != nil {
		updateBatchErrorTask(err, &taskInfo)
		logger.Errorln(err.Error())
		return err
	}
	filein, err := ioutil.ReadFile(filePath)
	if err != nil {
		updateBatchErrorTask(err, &taskInfo)
		logger.Errorln(err.Error())
		return err
	}
	r := csv.NewReader(strings.NewReader(string(filein)))
	records, err := r.ReadAll()
	if err != nil {
		updateBatchErrorTask(err, &taskInfo)
		logger.Errorln(err.Error())
		return err
	}
	titles := records[0]
	titlesMap := make(map[string]int)
	for i := 0; i < len(titles); i++ {
		titlesMap[titles[i]] = i
	}
	logger.Infoln(records)
	for i := 1; i < len(records); i++ {
		line := records[i]
		batchNode := &nodeentity.NodeOfBatchInfo{
			BatchID:   int64(batchID),
			ClusterId: taskInfo.ClusterId,
			GroupId:   taskInfo.GroupId,
			UserId:    taskInfo.UserId,
			NodeName:  line[titlesMap["name"]],
			NodeType:  line[titlesMap["nodeType"]],
			Arch:      line[titlesMap["arch"]],
			Address:   line[titlesMap["ip"]],
			Port:      line[titlesMap["port"]],
			Sudoer:    line[titlesMap["sudoer"]],
			Password:  line[titlesMap["password"]],
			IsConfirm: false,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		}
		err := nodeentity.CreateNodeOfBatch(batchNode)
		if err != nil {
			taskInfo.Status = "error"
			taskInfo.ErrMsg = err.Error()
			logger.Errorln(err.Error())
			break
		}
	}

	err = nodeentity.FinishBatchTask(batchID)
	if err != nil {
		updateBatchErrorTask(err, &taskInfo)
		logger.Errorln(err.Error())
		return err
	}
	return nil
}

func updateBatchErrorTask(err error, taskInfo *nodeentity.BatchTaskRecord) error {
	taskInfo.Status = "error"
	taskInfo.ErrMsg = err.Error()
	err = nodeentity.UpdateBatchTask(taskInfo)
	return err
}

func AddOneRecord(userInfo proto.ApulisHeader, content *nodemodule.CreateNodeOfBatchReq) error {
	batchNode := nodeentity.NodeOfBatchInfo{
		BatchID:   int64(content.ID),
		ClusterId: userInfo.ClusterId,
		GroupId:   userInfo.GroupId,
		UserId:    userInfo.UserId,
		NodeName:  content.Name,
		NodeType:  content.NodeType,
		Arch:      content.Arch,
		Address:   content.Address,
		Port:      content.Port,
		Sudoer:    content.Sudoer,
		Password:  content.Password,
		IsConfirm: false,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}

	err := nodeentity.CreateNodeOfBatch(&batchNode)
	if err != nil {
		return err
	}
	return nil
}

func ListBatchList(userInfo *proto.ApulisHeader, batchID int, pageSize int, pageNum int) (*nodeentity.BatchTaskRecord, *[]nodeentity.NodeOfBatchInfo, error) {
	var batchInfos []nodeentity.NodeOfBatchInfo
	var batchTask nodeentity.BatchTaskRecord

	res := apulisdb.Db.
		Where("BatchID = ? and ClusterId = ? and GroupId = ? and UserId = ? and isConfirm = ?", batchID, userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, false).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&batchInfos)

	if res.Error != nil {
		return &batchTask, &batchInfos, res.Error
	}

	batchTask, err := nodeentity.GetBatchTaskByID(batchID)
	if err != nil {
		return &batchTask, &batchInfos, err
	}

	return &batchTask, &batchInfos, nil
}

func UpdateBatch(userInfo proto.ApulisHeader, id int) error {
	var batchInfos []nodeentity.NodeOfBatchInfo

	res := apulisdb.Db.
		Where("ID = ? and ClusterId = ? and GroupId = ? and UserId = ? and isConfirm = ?", id, userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, false).
		Find(&batchInfos)
	if res.Error != nil {
		return res.Error
	}

	var err error
	clu, err := cluster.GetCluster(userInfo.ClusterId)
	if err != nil {
		logger.Errorln("CreateEdgeNode, can`t find cluster ", userInfo.ClusterId)
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
		batchInfo.NodeID = node.ID
		nodeentity.ConfirmNodesBatch(&batchInfo)
	}

	return nil
}

func DeleteNodeOfBatch(userInfo proto.ApulisHeader, req *nodemodule.DeleteNodeOfBatchReq) error {
	var nodeOfBatchInfo nodeentity.NodeOfBatchInfo

	res := apulisdb.Db.
		Where("ClusterId = ? and GroupId = ? and UserId = ? and ID = ?", userInfo.ClusterId, userInfo.GroupId, userInfo.UserId, req.ID).
		First(&nodeOfBatchInfo)
	if res.Error != nil {
		return res.Error
	}

	return nodeentity.DeleteNodeOfBatch(&nodeOfBatchInfo)
}

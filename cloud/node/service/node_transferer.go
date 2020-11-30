// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	"context"
	"fmt"
	apulisdb "github.com/apulis/ApulisEdge/cloud/database"
	constants "github.com/apulis/ApulisEdge/cloud/node"
	"time"
)

// CreateNodeCheckLoop transferer of edge node status
func CreateNodeCheckLoop(ctx context.Context, interval int32) {
	duration := time.Duration(interval) * time.Second
	checkTicker := time.NewTimer(duration)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Infof("CreateNodeCheckLoop was terminated")
			return
		case <-checkTicker.C:
			NodeChecker()
			checkTicker.Reset(duration)
		}
	}
}

func NodeChecker() {
	//var nodeInfo nodeentity.NodeBasicInfo

	queryStr := fmt.Sprintf("status = '%s'", constants.StatusNotInstalled)
	rows, _ := apulisdb.Db.Table("node_basic_infos").Where(queryStr).Select("name").Rows()
	//logger.Infof("NodeChecker find %d node not installed", result.RowsAffected)

	//rows, _ := result.Rows()
	//defer rows.Close()

	//cols, _ := rows.Columns()
	//logger.Info(cols)

	for rows.Next() {
		var n string
		rows.Scan(&n)
		logger.Info(n)
	}
}

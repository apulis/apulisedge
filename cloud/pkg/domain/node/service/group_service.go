// Copyright 2020 Apulis Technology Inc. All rights reserved.

package nodeservice

import (
	nodemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/node"
	nodeentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/node/entity"
	proto "github.com/apulis/ApulisEdge/cloud/pkg/protocol"
)

func CreateNodeGroup(userInfo proto.ApulisHeader, req *nodemodule.CreateEdgeNodeReq) (*nodeentity.NodeGroupInfo, error) {
	return nil, nil
}

func ListNodeGroup(userInfo proto.ApulisHeader, req *nodemodule.ListEdgeNodesReq) (*[]nodeentity.NodeGroupInfo, int, error) {
	return nil, 0, nil
}

func DescribeNodeGroup(userInfo proto.ApulisHeader, req *nodemodule.DescribeEdgeNodesReq) (*nodeentity.NodeBasicInfo, error) {
	return nil, nil
}

// delete edge group
func DeleteNodeGroup(userInfo proto.ApulisHeader, req *nodemodule.DeleteEdgeNodeReq) error {
	return nil
}

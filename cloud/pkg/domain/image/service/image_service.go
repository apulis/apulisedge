// Copyright 2020 Apulis Technology Inc. All rights reserved.

package imageservice

import (
	apulisdb "github.com/apulis/ApulisEdge/cloud/pkg/database"
	imagemodule "github.com/apulis/ApulisEdge/cloud/pkg/domain/image"
	imageentity "github.com/apulis/ApulisEdge/cloud/pkg/domain/image/entity"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
)

var logger = loggers.LogInstance()

func ListContainerImage(req *imagemodule.ListContainerImageReq) ([]imageentity.UserContainerImageInfo, int, error) {
	var imageInfos []imageentity.UserContainerImageInfo

	total := 0
	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	res := apulisdb.Db.Offset(offset).Limit(limit).
		Where("ClusterId = ? and GroupId = ? and UserId = ?", req.ClusterId, req.GroupId, req.UserId).
		Group("ImageName").
		Group("OrgName").
		Select("ImageName, OrgName", "UpdateAt").
		Find(&imageInfos)

	if res.Error != nil {
		return imageInfos, total, res.Error
	}

	return imageInfos, int(res.RowsAffected), nil
}

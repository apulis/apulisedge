// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"github.com/gin-gonic/gin"
)

// Get user info from context
func GetUserInfo(c *gin.Context) (int64, int64, int64, error) {
	clusterId, err := getClusterId(c)
	if err != nil {
		return InvalidClusterId, InvalidGroupId, InvalidUserId, ErrInvalidClusterId
	}

	groupId, err := getGroupId(c)
	if err != nil {
		return InvalidClusterId, InvalidGroupId, InvalidUserId, ErrInvalidGroupId
	}

	userId, err := getUserId(c)
	if err != nil {
		return InvalidClusterId, InvalidGroupId, InvalidUserId, ErrInvalidUserId
	}

	return clusterId, groupId, userId, nil
}

func getClusterId(c *gin.Context) (int64, error) {
	clusterId, exists := c.Get("clusterId")
	if !exists {
		return InvalidClusterId, ErrInvalidClusterId
	}

	return clusterId.(int64), nil
}

func getGroupId(c *gin.Context) (int64, error) {
	groupId, exists := c.Get("groupId")
	if !exists {
		return InvalidGroupId, ErrInvalidGroupId
	}

	return groupId.(int64), nil
}

func getUserId(c *gin.Context) (int64, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return InvalidUserId, ErrInvalidUserId
	}

	return userId.(int64), nil
}

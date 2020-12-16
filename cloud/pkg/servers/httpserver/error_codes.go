// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import "errors"

const (
	SUCCESS_CODE = 0

	NOT_FOUND_ERROR_CODE = 10001
	UNKNOWN_ERROR_CODE   = 10002
	SERVER_ERROR_CODE    = 10003

	// Request error codes
	PARAMETER_ERROR_CODE = 20001
	AUTH_ERROR_CODE      = 20002

	// APP error codes
	APP_ERROR_CODE              = 30000
	FILETYPE_NOT_SUPPORTED_CODE = 30001
	SAVE_FILE_ERROR_CODE        = 30002
	EXTRACT_FILE_ERROR_CODE     = 30003
	REMOVE_FILE_ERROR_CODE      = 30004
	FILEPATH_NOT_EXISTS_CODE    = 30005
	FILE_OVERSIZE_CODE          = 30006
	FILEPATH_NOT_VALID_CODE     = 30007
	COMPRESS_PATH_ERROR_CODE    = 30008

	// application
	APP_STILL_HAVE_DEPLOY = 30101

	REMOTE_SERVE_ERROR_CODE = 40000
)

// user manage
const (
	InvalidClusterId int64 = -1
	InvalidGroupId   int64 = -1
	InvalidUserId    int64 = -1
)

var (
	ErrInvalidClusterId = errors.New("Invalid cluster id")
	ErrInvalidGroupId   = errors.New("Invalid group id")
	ErrInvalidUserId    = errors.New("Invalid user id")
	ErrInvalidUserInfo  = errors.New("Invalid user info")
	ErrOrgNameNeeded    = errors.New("Need org name")
)

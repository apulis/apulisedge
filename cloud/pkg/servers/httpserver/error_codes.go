// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

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

	// dataset
	//上传大文件template目录满
	UPLOAD_TEMPDIR_FULL_CODE = 30101
	//无法删除正在使用的数据集
	DATASET_IS_STILL_USE_CODE = 30102
	//已经存在同名的数据集
	DATASET_IS_EXISTED = 30103

	REMOTE_SERVE_ERROR_CODE = 40000
)
// Copyright 2020 Apulis Technology Inc. All rights reserved.

package database

import (
	"database/sql"
	"fmt"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log"
)

var Db *gorm.DB
var logger = loggers.LogInstance()

func InitDatabase(config *configs.EdgeCloudConfig) {
	dbConf := config.Db

	sqlDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database))
	if err != nil {
		panic(err)
	}

	_, err = sqlDb.Exec("CREATE DATABASE IF NOT EXISTS " + dbConf.Database)
	if err != nil {
		panic(err)
	}

	var lvl gormlogger.LogLevel
	if config.DebugModel {
		lvl = gormlogger.Info
	} else {
		lvl = gormlogger.Warn
	}

	newLogger := gormlogger.New(
		log.New(logger.Out, "\r\n", log.LstdFlags),
		gormlogger.Config{
			LogLevel: lvl,
		},
	)

	Db, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDb,
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	logger.Info("DB connected success")
	sqlDb.SetMaxOpenConns(dbConf.MaxOpenConns)
	sqlDb.SetMaxIdleConns(dbConf.MaxIdleConns)
}

func CreateTableIfNotExists(modelType interface{}) error {
	var err error
	if err = Db.
		Set("gorm:table_options", "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8").
		AutoMigrate(modelType); err != nil {
		logger.Errorf("AutoMigrate failed! table = %s", modelType)
	}
	return err
}

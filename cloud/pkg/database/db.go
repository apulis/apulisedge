// Copyright 2020 Apulis Technology Inc. All rights reserved.

package database

import (
	"database/sql"
	"fmt"
	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
	"github.com/apulis/ApulisEdge/cloud/pkg/loggers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var logger = loggers.LogInstance()

func InitDatabase(config *configs.EdgeCloudConfig) {
	dbConf := config.Db

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port))

	defer db.Close()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbConf.Database)
	if err != nil {
		panic(err)
	}

	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger.Info("DB connected success")
	//Db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	//Db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
}

func CreateTableIfNotExists(modelType interface{}) error {
	var err error
	if err = Db.AutoMigrate(modelType); err != nil {
		logger.Errorf("AutoMigrate failed! table = %s", modelType)
	}
	return err
}

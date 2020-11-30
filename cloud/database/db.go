// Copyright 2020 Apulis Technology Inc. All rights reserved.

package database

import (
	"database/sql"
	"fmt"
	"github.com/apulis/ApulisEdge/cloud/configs"
	"github.com/apulis/ApulisEdge/cloud/loggers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
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

	Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database))
	if err != nil {
		panic(err)
	}

	logger.Info("DB connected success")
	Db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	Db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
}

func CreateTableIfNotExists(modelType interface{}) {
	val := reflect.Indirect(reflect.ValueOf(modelType))
	modelName := val.Type().Name()

	hasTable := Db.HasTable(modelType)
	if !hasTable {
		logger.Info(fmt.Sprintf("Table of %s not exists, create it.", modelName))
		Db.CreateTable(modelType)
	} else {
		logger.Info(fmt.Sprintf("Table of %s already exists.", modelName))
	}
}

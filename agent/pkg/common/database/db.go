package database

import (
	"os"
	"path"
	"reflect"

	"github.com/apulis/ApulisEdge/agent/pkg/common/config"
	"github.com/apulis/ApulisEdge/agent/pkg/common/loggers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var logger = loggers.LogInstance()

// Db is the database connection
var Db *gorm.DB

func InitDatabase() error {
	// create database dir
	_, err := os.Stat(config.AppConfig.Database.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(config.AppConfig.Database.Dir, 0755)
			logger.Info("Create Database directory: %s", config.AppConfig.Database.Dir)
		} else {
			logger.Panicln("Fatal database directory error: %s", err)
		}
	}

	// connect database
	if config.AppConfig.Database.Type == "sqlite3" {
		sqlite_db_file_path := path.Join(config.AppConfig.Database.Dir, config.SQLITE_DB_FILE_NAME)
		if Db, err = gorm.Open(config.AppConfig.Database.Type, sqlite_db_file_path); err != nil {
			logger.Panicln("Fatal open database: %s", err)
		}
	} else {
		logger.Panicln("Fatal support database type: ", config.AppConfig.Database.Type)
	}
	logger.Infoln("Database connected succeed")

	createDatabaseIfNotExists()

	return nil
}

func createDatabaseIfNotExists() {
}

func CreateTableIfNotExists(modelType interface{}) {
	val := reflect.Indirect(reflect.ValueOf(modelType))
	modelName := val.Type().Name()

	hasTable := Db.HasTable(modelType)
	if !hasTable {
		logger.Infoln("Table of ", modelName, "not exists, create it.")
		Db.CreateTable(modelType)
	} else {
		logger.Infoln("Table of ", modelName, "already exists.")
	}
}

func CloseDatabase() {
	defer Db.Close()
}

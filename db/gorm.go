package db

import (
	"go-kurs-bca/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	// Register Gorm Mysql Driver
	_ "gorm.io/driver/mysql"
	// Register Go Sql Driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	appConfig = config.Config.App
	dbConfig  = config.Config.DB
	mysqlConn *gorm.DB
	err       error
)

func init() {
	if dbConfig.Driver == "mysql" {
		setupMySQLConn()
	}
}

func setupMySQLConn() {
	var logLevel logger.LogLevel = logger.Silent

	if appConfig.Env == "local" {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
			Colorful:      true,        // Disable/Enable color
		},
	)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	mysqlConn, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
}

func MySQLConn() *gorm.DB {
	return mysqlConn
}

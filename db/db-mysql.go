package db

import (
	"goTestProject/config"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbMySQL = gorm.DB{}

func init() {
	initMySQLDB("gotestproject")
}

func initMySQLDB(dbName string) {
	var datetimePrecision = 2
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.Conf.Common.CommonMysql.MysqlAddress, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                                         // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,                                        // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision,                          // default datetime precision
		DontSupportRenameIndex:    true,                                        // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                        // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                       // smart configure based on used version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Error("connect db fail:%s", err.Error())
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbMySQL = *db
}

func GetMySQLDb() (db *gorm.DB) {
	return &dbMySQL
}

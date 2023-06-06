package db

import (
	"crazy_server/pkg/common/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type agentMysqlDB struct {
	//sync.RWMutex
	db *gorm.DB
}

type AgentWriter struct{}

func (w AgentWriter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func initAgentMysqlDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.AgentMysql.DBUserName, config.Config.AgentMysql.DBPassword, config.Config.AgentMysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), nil)
		if err1 != nil {
			panic(err1.Error() + " open failed " + dsn)
		}
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.AgentMysql.DBUserName, config.Config.AgentMysql.DBPassword, config.Config.AgentMysql.DBAddress[0], config.Config.AgentMysql.DBDatabaseName)
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(config.Config.AgentMysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(config.Config.AgentMysql.LogLevel),                       // Log level
			IgnoreRecordNotFoundError: true,                                                                     // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                                     // Disable color
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err.Error() + " Open failed " + dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error() + " db.DB() failed ")
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.AgentMysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.AgentMysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.AgentMysql.DBMaxIdleConns)

	DB.AgentMysqlDB.db = db
}

func (m *agentMysqlDB) DefaultGormDB() *gorm.DB {
	return DB.AgentMysqlDB.db
}

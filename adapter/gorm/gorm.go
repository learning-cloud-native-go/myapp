package gorm

import (
	"fmt"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"myapp/config"
)

func New(conf *config.DbConf) (*gorm.DB, error) {
	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		DBName:               conf.DbName,
		User:                 conf.Username,
		Passwd:               conf.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	var logLevel logger.LogLevel
	if conf.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
}

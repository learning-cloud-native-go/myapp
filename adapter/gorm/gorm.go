package gorm

import (
	"fmt"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"myapp/config"
)

func New(conf *config.Conf) (*gorm.DB, error) {
	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Db.Host, conf.Db.Port),
		DBName:               conf.Db.DbName,
		User:                 conf.Db.Username,
		Passwd:               conf.Db.Password,
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

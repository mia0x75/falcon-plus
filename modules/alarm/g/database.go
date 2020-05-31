package g

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	db *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB() (err error) {
	cfg := Config()
	db, err := gorm.Open("mysql", cfg.Database.Addr)
	if err != nil {
		return fmt.Errorf("connect to dashboard: %s", err.Error())
	}
	db.SingularTable(true)
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	}

	db.DB().SetMaxIdleConns(cfg.Database.MaxIdle)
	db.DB().SetMaxOpenConns(cfg.Database.MaxConnections)
	db.DB().SetConnMaxLifetime(time.Duration(cfg.Database.WaitTimeout) * time.Second)

	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	return
}

// Con 获取链接
func Con() *gorm.DB {
	return db
}

// SetLogLevel 设置日志级别
func SetLogLevel(loggerlevel bool) {
	db.LogMode(loggerlevel)
}

// CloseDB 关闭链接
func CloseDB() (err error) {
	err = db.Close()
	if err != nil {
		return
	}
	return
}

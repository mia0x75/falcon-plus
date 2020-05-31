package g

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	db *gorm.DB
)

func Con() *gorm.DB {
	return db
}

func SetLogLevel(loggerlevel bool) {
	db.LogMode(loggerlevel)
}

func InitDB() (err error) {
	var d *sql.DB
	dashd, err := gorm.Open("mysql", Config().Database.Addr)
	dashd.Dialect().SetDB(d)
	if err != nil {
		return fmt.Errorf("connect to dashboard: %s", err.Error())
	}
	dashd.SingularTable(true)
	dashd.DB().SetMaxIdleConns(Config().Database.MaxIdle)
	dashd.DB().SetMaxOpenConns(Config().Database.MaxConnections)
	dashd.DB().SetConnMaxLifetime(time.Duration(Config().Database.WaitTimeout) * time.Second)
	err = dashd.DB().Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	db = dashd
	return
}

func CloseDB() (err error) {
	err = db.Close()
	if err != nil {
		return
	}
	return
}

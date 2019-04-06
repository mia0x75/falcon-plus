package service

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

var DB *sql.DB

func InitDB() (err error) {
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	}

	DB.SetMaxOpenConns(g.Config().Database.MaxConnections)
	DB.SetMaxIdleConns(g.Config().Database.MaxIdle)
	DB.SetConnMaxLifetime(time.Duration(g.Config().Database.WaitTimeout) * time.Second)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	return
}

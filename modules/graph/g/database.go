package g

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func InitDB() (err error) {
	DB, err = sql.Open("mysql", Config().Database.Addr)
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	}

	DB.SetMaxOpenConns(Config().Database.MaxConnections)
	DB.SetMaxIdleConns(Config().Database.MaxIdle)
	DB.SetConnMaxLifetime(time.Duration(Config().Database.WaitTimeout) * time.Second)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	return
}

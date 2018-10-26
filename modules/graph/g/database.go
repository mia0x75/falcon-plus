package g

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() (err error) {
	DB, err = sql.Open("mysql", Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(Config().Database.MaxConnections)
	DB.SetMaxIdleConns(Config().Database.MaxIdle)
	DB.SetConnMaxLifetime(time.Duration(Config().Database.WaitTimeout) * time.Second)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
	return
}

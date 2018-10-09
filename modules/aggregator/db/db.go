package db

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxIdleConns(g.Config().Database.MaxIdle)
	DB.SetMaxOpenConns(g.Config().Database.MaxConnections)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}

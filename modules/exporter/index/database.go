package index

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = GetDbConn()
	if err != nil {
		log.Fatalln("open db fail:", err)
	} else {
		log.Println("index:InitDB ok")
	}
}

func GetDbConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", g.Config().Index.Addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(g.Config().Index.MaxIdle)
	db.SetMaxOpenConns(g.Config().Index.MaxConnections)
	db.SetConnMaxLifetime(time.Duration(g.Config().Index.WaitTimeout) * time.Second)

	err = db.Ping()
	if err != nil {
		db.Close()
	}
	return
}

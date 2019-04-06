package index

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = GetDbConn()
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	} else {
		log.Infof("[I] index:InitDB ok")
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

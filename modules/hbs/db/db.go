package db

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(g.Config().Database.MaxConns)
	DB.SetMaxIdleConns(g.Config().Database.MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}

package index

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // TODO:
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/exporter/g"
)

var (
	db *sql.DB
)

// Con 获取链接
func Con() *sql.DB {
	return db
}

// Close 关闭链接
func Close() {
	db.Close()
}

// InitDB 初始化数据库连接
func InitDB() (err error) {
	db, err = sql.Open("mysql", g.Config().Index.Addr)
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	}

	db.SetMaxIdleConns(g.Config().Index.MaxIdle)
	db.SetMaxOpenConns(g.Config().Index.MaxConnections)
	db.SetConnMaxLifetime(time.Duration(g.Config().Index.WaitTimeout) * time.Second)

	err = db.Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	return
}

package g

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DBPool struct {
	Falcon    *gorm.DB
	Graph     *gorm.DB
	Uic       *gorm.DB
	Dashboard *gorm.DB
	Alarm     *gorm.DB
}

var (
	dbp DBPool
)

func Con() DBPool {
	return dbp
}

func SetLogLevel(loggerlevel bool) {
	dbp.Uic.LogMode(loggerlevel)
	dbp.Graph.LogMode(loggerlevel)
	dbp.Falcon.LogMode(loggerlevel)
	dbp.Dashboard.LogMode(loggerlevel)
	dbp.Alarm.LogMode(loggerlevel)
}

func InitDB() (err error) {
	var p *sql.DB
	portald, err := gorm.Open("mysql", Config().Databases.Portal.Addr)
	portald.Dialect().SetDB(p)
	if err != nil {
		return fmt.Errorf("connect to portal: %s", err.Error())
	}
	portald.SingularTable(true)
	portald.DB().SetMaxIdleConns(Config().Databases.Portal.MaxIdle)
	portald.DB().SetMaxOpenConns(Config().Databases.Portal.MaxConnections)
	portald.DB().SetConnMaxLifetime(time.Duration(Config().Databases.Portal.WaitTimeout) * time.Second)
	dbp.Falcon = portald

	var g *sql.DB
	graphd, err := gorm.Open("mysql", Config().Databases.Graph.Addr)
	graphd.Dialect().SetDB(g)
	if err != nil {
		return fmt.Errorf("connect to graph: %s", err.Error())
	}
	graphd.SingularTable(true)
	graphd.DB().SetMaxIdleConns(Config().Databases.Graph.MaxIdle)
	graphd.DB().SetMaxOpenConns(Config().Databases.Graph.MaxConnections)
	graphd.DB().SetConnMaxLifetime(time.Duration(Config().Databases.Graph.WaitTimeout) * time.Second)
	dbp.Graph = graphd

	var u *sql.DB
	uicd, err := gorm.Open("mysql", Config().Databases.Uic.Addr)
	uicd.Dialect().SetDB(u)
	if err != nil {
		return fmt.Errorf("connect to uic: %s", err.Error())
	}
	uicd.SingularTable(true)
	uicd.DB().SetMaxIdleConns(Config().Databases.Uic.MaxIdle)
	uicd.DB().SetMaxOpenConns(Config().Databases.Uic.MaxConnections)
	uicd.DB().SetConnMaxLifetime(time.Duration(Config().Databases.Uic.WaitTimeout) * time.Second)
	dbp.Uic = uicd

	var d *sql.DB
	dashd, err := gorm.Open("mysql", Config().Databases.Dashboard.Addr)
	dashd.Dialect().SetDB(d)
	if err != nil {
		return fmt.Errorf("connect to dashboard: %s", err.Error())
	}
	dashd.SingularTable(true)
	dashd.DB().SetMaxIdleConns(Config().Databases.Dashboard.MaxIdle)
	dashd.DB().SetMaxOpenConns(Config().Databases.Dashboard.MaxConnections)
	dashd.DB().SetConnMaxLifetime(time.Duration(Config().Databases.Dashboard.WaitTimeout) * time.Second)
	dbp.Dashboard = dashd

	var alm *sql.DB
	almd, err := gorm.Open("mysql", Config().Databases.Alarms.Addr)
	almd.Dialect().SetDB(alm)
	if err != nil {
		return fmt.Errorf("connect to alarms: %s", err.Error())
	}
	almd.SingularTable(true)
	almd.DB().SetMaxIdleConns(Config().Databases.Alarms.MaxIdle)
	almd.DB().SetMaxOpenConns(Config().Databases.Alarms.MaxConnections)
	almd.DB().SetConnMaxLifetime(time.Duration(Config().Databases.Alarms.WaitTimeout) * time.Second)
	dbp.Alarm = almd

	return
}

func CloseDB() (err error) {
	err = dbp.Falcon.Close()
	if err != nil {
		return
	}
	err = dbp.Graph.Close()
	if err != nil {
		return
	}
	err = dbp.Uic.Close()
	if err != nil {
		return
	}
	err = dbp.Dashboard.Close()
	if err != nil {
		return
	}
	err = dbp.Alarm.Close()
	if err != nil {
		return
	}
	return
}

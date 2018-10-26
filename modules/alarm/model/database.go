package model

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func InitDB() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Database.Addr, config.Database.MaxIdle, config.Database.MaxConnections)
	// register model
	orm.RegisterModel(new(event.Events), new(event.EventCases))
	if db, err := orm.GetDB(); err != nil {
		log.Fatalln("open db fail:", err)
	} else {
		db.SetConnMaxLifetime(time.Duration(config.Database.WaitTimeout) * time.Second)
	}
	if g.IsDebug() {
		orm.Debug = true
	}
}

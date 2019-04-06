package model

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
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
		log.Fatalf("[F] open db fail: %v", err)
	} else {
		db.SetConnMaxLifetime(time.Duration(config.Database.WaitTimeout) * time.Second)

		err = db.Ping()
		if err != nil {
			log.Fatalf("[F] ping db fail: %v", err)
		}
	}

	if cutils.IsDebug() {
		orm.Debug = true
	}
}

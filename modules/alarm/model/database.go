package model

import (
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
	if g.IsDebug() {
		orm.Debug = true
	}
}

package model

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Portal.Addr, config.Portal.Idle, config.Portal.Max)
	// register model
	orm.RegisterModel(new(event.Events), new(event.EventCases))
	if config.LogLevel == "debug" {
		orm.Debug = true
	}
}

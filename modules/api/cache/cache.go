package cache

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/api/g"
)

var db *gorm.DB

func Init() {
	db = g.Con()

	log.Info("[I] cache begin")

	log.Info("[I] #01 UsersMap...")
	UsersMap.Init()

	log.Info("[I] #02 TeamsMap...")
	TeamsMap.Init()

	log.Info("[I] #03 HostsMap...")
	HostsMap.Init()

	log.Info("[I] #04 ExpressionsMap...")
	ExpressionsMap.Init()

	log.Info("[I] #05 ActionsMap...")
	ActionsMap.Init()

	log.Info("[I] #06 ClustersMap...")
	ClustersMap.Init()

	log.Info("[I] #07 GroupsMap...")
	GroupsMap.Init()

	log.Info("[I] #08 StrategiesMap...")
	StrategiesMap.Init()

	log.Info("[I] #09 TemplatesMap...")
	TemplatesMap.Init()

	log.Info("[I] #10 EdgesMap...")
	EdgesMap.Init()

	log.Info("[I] cache done")

	go LoopInit()
}

func LoopInit() {
	d := time.Duration(1) * time.Minute
	for range time.Tick(d) {
		UsersMap.Init()
		TeamsMap.Init()
		HostsMap.Init()
		ExpressionsMap.Init()
		ActionsMap.Init()
		ClustersMap.Init()
		GroupsMap.Init()
		StrategiesMap.Init()
		TemplatesMap.Init()
	}
}

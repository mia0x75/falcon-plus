package cache

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

func Init() {
	log.Println("cache begin")

	log.Println("#0 GroupPlugins...")
	GroupPlugins.Init()

	log.Println("#1 GroupTemplates...")
	GroupTemplates.Init()

	log.Println("#2 HostGroupsMap...")
	HostGroupsMap.Init()

	log.Println("#3 HostMap...")
	HostMap.Init()

	log.Println("#4 TemplateCache...")
	TemplateCache.Init()

	log.Println("#5 Strategies...")
	Strategies.Init(TemplateCache.GetMap())

	log.Println("#6 HostTemplateIds...")
	HostTemplateIds.Init()

	log.Println("#7 ExpressionCache...")
	ExpressionCache.Init()

	log.Println("#8 MonitoredHosts...")
	MonitoredHosts.Init()

	log.Println("#9 AgentsInfo...")
	Agents.Init()

	log.Println("cache done")

	go LoopInit()
}

func LoopInit() {
	d := time.Duration(1) * time.Minute
	for range time.Tick(d) {
		GroupPlugins.Init()
		GroupTemplates.Init()
		HostGroupsMap.Init()
		HostMap.Init()
		TemplateCache.Init()
		Strategies.Init(TemplateCache.GetMap())
		HostTemplateIds.Init()
		ExpressionCache.Init()
		MonitoredHosts.Init()
		Agents.Init()
	}
}

package cache

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.Info("[I] cache begin")

	log.Info("[I] #0 GroupPlugins...")
	GroupPlugins.Init()

	log.Info("[I] #1 GroupTemplates...")
	GroupTemplates.Init()

	log.Info("[I] #2 HostGroupsMap...")
	HostGroupsMap.Init()

	log.Info("[I] #3 HostMap...")
	HostMap.Init()

	log.Info("[I] #4 TemplateCache...")
	TemplateCache.Init()

	log.Info("[I] #5 Strategies...")
	Strategies.Init(TemplateCache.GetMap())

	log.Info("[I] #6 HostTemplateIds...")
	HostTemplateIds.Init()

	log.Info("[I] #7 ExpressionCache...")
	ExpressionCache.Init()

	log.Info("[I] #8 MonitoredHosts...")
	MonitoredHosts.Init()

	log.Info("[I] #9 AgentsInfo...")
	Agents.Init()

	log.Info("[I] cache done")

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

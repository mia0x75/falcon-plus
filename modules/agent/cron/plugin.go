package cron

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/plugins"
)

func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if len(g.Config().Heartbeat.Addrs) == 0 {
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {
	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := cmodel.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp cmodel.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp

		log.Debugln(&resp)

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
		}

		desiredAll := make(map[string]*plugins.Plugin)

		for _, p := range pluginDirs {
			underOneDir := plugins.ListPlugins(strings.Trim(p, "/"))
			for k, v := range underOneDir {
				desiredAll[k] = v
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)
	}
}

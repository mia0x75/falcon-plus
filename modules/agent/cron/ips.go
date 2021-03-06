package cron

import (
	"time"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// SyncTrustableIps TODO:
func SyncTrustableIps() {
	if len(g.Config().Heartbeat.Addrs) > 0 {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {
	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", cm.NullRPCRequest{}, &ips)
		if err != nil {
			log.Errorf("[E] call Agent.TrustableIps fail: %v", err)
			continue
		}

		hbs.CacheTrustableIps(ips)
	}
}

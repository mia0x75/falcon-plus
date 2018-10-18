package cron

import (
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func SyncTrustableIps() {
	if len(g.Config().Heartbeat.Addrs) > 0 {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {
	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", cmodel.NullRpcRequest{}, &ips)
		if err != nil {
			log.Println("ERROR: call Agent.TrustableIps fail", err)
			continue
		}

		g.SetTrustableIps(ips)
	}
}

package cron

import (
	"log"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func SyncTrustableIps() {
	if g.Config().Heartbeat.Enabled && len(g.Config().Heartbeat.Addrs) > 0 {
		go syncTrustableIps()
	}
}

func syncTrustableIps() {
	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		var ips string
		err := g.HbsClient.Call("Agent.TrustableIps", model.NullRpcRequest{}, &ips)
		if err != nil {
			log.Println("ERROR: call Agent.TrustableIps fail", err)
			continue
		}

		g.SetTrustableIps(ips)
	}
}

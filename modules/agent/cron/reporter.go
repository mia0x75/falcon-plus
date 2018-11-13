package cron

import (
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func ReportAgentStatus() {
	if len(g.Config().Heartbeat.Addrs) > 0 {
		go reportAgentStatus()
	}
}

func reportAgentStatus() {
	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		hostname, err := g.Hostname()
		if err != nil {
			log.Errorf("[E] get hostname error: %v", err)
			continue
		}

		req := cmodel.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: g.GetCurrPluginVersion(),
			// TODO: Add system information to support inventory management
		}

		var resp cmodel.SimpleRpcResponse
		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			log.Errorf("[E] call Agent.ReportStatus fail: %v Request: %v Response: %v\n", err, req, resp)
		} else {
			log.Debugf("[D] call Agent.ReportStatus success. Request: %v Response: %v\n", req, resp)
		}
	}
}

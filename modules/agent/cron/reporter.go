package cron

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// ReportAgentStatus 向心跳服务器上报Agent信息
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
			AgentVersion:  fmt.Sprintf("%s@%s", g.Version, g.Commit),
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

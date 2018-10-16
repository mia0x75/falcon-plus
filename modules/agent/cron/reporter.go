package cron

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
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
			hostname = fmt.Sprintf("error:%s", err.Error())
		}

		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: g.GetCurrPluginVersion(),
			// TODO: Add system information to support inventory management
		}

		var resp model.SimpleRpcResponse
		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			log.Errorf("call Agent.ReportStatus fail: %v Request: %v Response: %v\n", err, req, resp)
		} else {
			log.Debugf("call Agent.ReportStatus success. Request: %v Response: %v\n", req, resp)
		}
	}
}

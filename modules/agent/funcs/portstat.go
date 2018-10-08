package funcs

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
)

func PortMetrics() (L []*model.MetricValue) {
	reportPorts := g.ReportPorts()
	sz := len(reportPorts)
	if sz == 0 {
		return
	}

	allTcpPorts, err := nux.TcpPorts()
	if err != nil {
		log.Println(err)
		return
	}

	allUdpPorts, err := nux.UdpPorts()
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", reportPorts[i])
		if slice.ContainsInt64(allTcpPorts, reportPorts[i]) || slice.ContainsInt64(allUdpPorts, reportPorts[i]) {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 1, tags))
		} else {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 0, tags))
		}
	}

	return
}

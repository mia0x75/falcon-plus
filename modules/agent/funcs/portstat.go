package funcs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// PortMetrics TODO:
func PortMetrics() (L []*cm.MetricValue) {
	ports := hbs.ReportPorts()
	sz := len(ports)
	if sz == 0 {
		return
	}

	allTCPPorts, err := nux.TcpPorts()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	allUDPPorts, err := nux.UdpPorts()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", ports[i])
		if slice.ContainsInt64(allTCPPorts, ports[i]) || slice.ContainsInt64(allUDPPorts, ports[i]) {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 1, tags))
		} else {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 0, tags))
		}
	}

	return
}

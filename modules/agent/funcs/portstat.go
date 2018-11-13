package funcs

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
)

func PortMetrics() (L []*cmodel.MetricValue) {
	ports := hbs.ReportPorts()
	sz := len(ports)
	if sz == 0 {
		return
	}

	allTcpPorts, err := nux.TcpPorts()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	allUdpPorts, err := nux.UdpPorts()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	for i := 0; i < sz; i++ {
		tags := fmt.Sprintf("port=%d", ports[i])
		if slice.ContainsInt64(allTcpPorts, ports[i]) || slice.ContainsInt64(allUdpPorts, ports[i]) {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 1, tags))
		} else {
			L = append(L, GaugeValue(g.NET_PORT_LISTEN, 0, tags))
		}
	}

	return
}

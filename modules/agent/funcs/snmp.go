package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// UDPMetrics TODO:
func UDPMetrics() []*cm.MetricValue {
	udp, err := nux.Snmp("Udp")
	if err != nil {
		log.Errorf("[E] read snmp fail: %v", err)
		return []*cm.MetricValue{}
	}

	count := len(udp)
	ret := make([]*cm.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.udp."+key, val)
		i++
	}

	return ret
}

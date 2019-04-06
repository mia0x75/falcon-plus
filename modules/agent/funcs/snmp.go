package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func UdpMetrics() []*cmodel.MetricValue {
	udp, err := nux.Snmp("Udp")
	if err != nil {
		log.Errorf("[E] read snmp fail: %v", err)
		return []*cmodel.MetricValue{}
	}

	count := len(udp)
	ret := make([]*cmodel.MetricValue, count)
	i := 0
	for key, val := range udp {
		ret[i] = CounterValue("snmp.udp."+key, val)
		i++
	}

	return ret
}

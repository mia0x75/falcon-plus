package funcs

import (
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
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

package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// SocketStatSummaryMetrics TODO:
func SocketStatSummaryMetrics() (L []*cm.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}

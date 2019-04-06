package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func SocketStatSummaryMetrics() (L []*cmodel.MetricValue) {
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

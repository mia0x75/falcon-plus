package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// LoadAvgMetrics TODO:
func LoadAvgMetrics() []*cm.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Errorf("[E] %v", err)
		return nil
	}

	return []*cm.MetricValue{
		GaugeValue("load.1min", load.Avg1min),
		GaugeValue("load.5min", load.Avg5min),
		GaugeValue("load.15min", load.Avg15min),
	}
}

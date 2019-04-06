package funcs

import (
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/nux"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func LoadAvgMetrics() []*cmodel.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Errorf("[E] %v", err)
		return nil
	}

	return []*cmodel.MetricValue{
		GaugeValue("load.1min", load.Avg1min),
		GaugeValue("load.5min", load.Avg5min),
		GaugeValue("load.15min", load.Avg15min),
	}
}

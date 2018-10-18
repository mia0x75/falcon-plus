package funcs

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func AgentMetrics() []*cmodel.MetricValue {
	return []*cmodel.MetricValue{GaugeValue("agent.alive", 1)}
}

package funcs

import (
	"strings"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *cmodel.MetricValue {
	mv := cmodel.MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}

	size := len(tags)

	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}

	return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *cmodel.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *cmodel.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}

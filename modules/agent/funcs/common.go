package funcs

import (
	"strings"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// NewMetricValue TODO:
func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *cm.MetricValue {
	mv := cm.MetricValue{
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

// GaugeValue TODO:
func GaugeValue(metric string, val interface{}, tags ...string) *cm.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

// CounterValue TODO:
func CounterValue(metric string, val interface{}, tags ...string) *cm.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}

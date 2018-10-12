package metric

import (
	"strings"

	metrics "github.com/mia0x75/go-metrics"
)

// self
func selfGauge(name string, value int64) {
	rr := gpSelf.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Gauge); ok {
			r.Update(value)
		}
		return
	}

	r := metrics.NewGauge()
	r.Update(value)
	if err := gpSelf.Register(name, r); isDuplicateMetricError(err) {
		r := gpSelf.Get(name).(metrics.Gauge)
		r.Update(value)
	}
}

func selfMeter(name string, value int64) {
	rr := gpSelf.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			r.Mark(value)
		}
		return
	}

	r := metrics.NewMeter()
	r.Mark(value)
	if err := gpSelf.Register(name, r); isDuplicateMetricError(err) {
		r := gpSelf.Get(name).(metrics.Meter)
		r.Mark(value)
	}
}

// internal
func isDuplicateMetricError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Index(err.Error(), "duplicate metric:") == 0
}

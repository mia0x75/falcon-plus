package metric

import (
	metrics "github.com/mia0x75/go-metrics"
)

// gauge
func Gauge(name string, value int64) {
	SetGaugeValue(name, float64(value))
}

func GaugeFloat64(name string, value float64) {
	SetGaugeValue(name, value)
}

func SetGaugeValue(name string, value float64) {
	rr := gpGaugeFloat64.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.GaugeFloat64); ok {
			r.Update(value)
		}
		return
	}

	r := metrics.NewGaugeFloat64()
	r.Update(value)
	if err := gpGaugeFloat64.Register(name, r); isDuplicateMetricError(err) {
		r := gpGaugeFloat64.Get(name).(metrics.GaugeFloat64)
		r.Update(value)
	}
}

func GetGaugeValue(name string) float64 {
	rr := gpGaugeFloat64.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.GaugeFloat64); ok {
			return r.Value()
		}
	}
	return 0.0
}

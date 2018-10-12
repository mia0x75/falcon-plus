package metric

import (
	metrics "github.com/mia0x75/go-metrics"
)

// meter
func Meter(name string, count int64) {
	SetMeterCount(name, count)
}

func SetMeterCount(name string, count int64) {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			r.Mark(count)
		}
		return
	}

	r := metrics.NewMeter()
	r.Mark(count)
	if err := gpMeter.Register(name, r); isDuplicateMetricError(err) {
		r := gpMeter.Get(name).(metrics.Meter)
		r.Mark(count)
	}
}

func GetMeterCount(name string) int64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Count()
		}
	}
	return 0
}

func GetMeterRateStep(name string) float64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.RateStep()
		}
	}
	return 0.0
}

func GetMeterRateMean(name string) float64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.RateMean()
		}
	}
	return 0.0
}

func GetMeterRate1(name string) float64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate1()
		}
	}
	return 0.0
}

func GetMeterRate5(name string) float64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate5()
		}
	}
	return 0.0
}

func GetMeterRate15(name string) float64 {
	rr := gpMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate15()
		}
	}
	return 0.0
}

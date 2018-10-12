package metric

import (
	metrics "github.com/mia0x75/go-metrics"
)

// senior
func Counter(name string, count int64) {
	SetCounterCount(name, count)
}

func SetCounterCount(name string, count int64) {
	rr := gpCounter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Counter); ok {
			r.Inc(count)
		}
		return
	}

	r := metrics.NewCounter()
	r.Inc(count)
	if err := gpCounter.Register(name, r); isDuplicateMetricError(err) {
		r := gpCounter.Get(name).(metrics.Counter)
		r.Inc(count)
	}
}

func GetCounterCount(name string) int64 {
	rr := gpCounter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Counter); ok {
			return r.Count()
		}
	}
	return 0
}

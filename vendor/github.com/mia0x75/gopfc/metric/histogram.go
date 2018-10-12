package metric

import (
	metrics "github.com/mia0x75/go-metrics"
)

// histogram
func Histogram(name string, count int64) {
	SetHistogramCount(name, count)
}

func SetHistogramCount(name string, count int64) {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			r.Update(count)
		}
		return
	}

	s := metrics.NewExpDecaySample(1028, 0.015)
	r := metrics.NewHistogram(s)
	r.Update(count)
	if err := gpHistogram.Register(name, r); isDuplicateMetricError(err) {
		r := gpHistogram.Get(name).(metrics.Histogram)
		r.Update(count)
	}
}

func GetHistogramCount(name string) int64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Count()
		}
	}
	return 0
}

func GetHistogramMax(name string) int64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Max()
		}
	}
	return 0
}

func GetHistogramMin(name string) int64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Min()
		}
	}
	return 0
}

func GetHistogramSum(name string) int64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Sum()
		}
	}
	return 0
}

func GetHistogramMean(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Mean()
		}
	}
	return 0.0
}

func GetHistogramStdDev(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.StdDev()
		}
	}
	return 0.0
}

func GetHistogram50th(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.5)
		}
	}
	return 0.0
}

func GetHistogram75th(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.75)
		}
	}
	return 0.0
}

func GetHistogram95th(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.95)
		}
	}
	return 0.0
}

func GetHistogram99th(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.99)
		}
	}
	return 0.0
}

func GetHistogram999th(name string) float64 {
	rr := gpHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.999)
		}
	}
	return 0.0
}

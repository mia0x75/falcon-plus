package funcs

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// FuncsAndInterval TODO:
type FuncsAndInterval struct {
	Fs       []func() []*cmodel.MetricValue
	Interval int
}

// Mappers TODO:
var Mappers []FuncsAndInterval

// BuildMappers TODO:
func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*cmodel.MetricValue{
				CPUMetrics,
				NetMetrics,
				KernelMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
				NetstatMetrics,
				ProcMetrics,
				UDPMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				PortMetrics,
				SocketStatSummaryMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				DuMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				URLMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				MySQLMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cmodel.MetricValue{
				RedisMetrics,
			},
			Interval: interval,
		},
	}
}

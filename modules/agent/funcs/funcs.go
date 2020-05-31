package funcs

import (
	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// FuncsAndInterval TODO:
type FuncsAndInterval struct {
	Fs       []func() []*cm.MetricValue
	Interval int
}

// Mappers TODO:
var Mappers []FuncsAndInterval

// BuildMappers TODO:
func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*cm.MetricValue{
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
			Fs: []func() []*cm.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cm.MetricValue{
				PortMetrics,
				SocketStatSummaryMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cm.MetricValue{
				DuMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cm.MetricValue{
				URLMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cm.MetricValue{
				MySQLMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*cm.MetricValue{
				RedisMetrics,
			},
			Interval: interval,
		},
	}
}

package funcs

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

type FuncsAndInterval struct {
	Fs       []func() []*cmodel.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*cmodel.MetricValue{
				CpuMetrics,
				NetMetrics,
				KernelMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
				NetstatMetrics,
				ProcMetrics,
				UdpMetrics,
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
				UrlMetrics,
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

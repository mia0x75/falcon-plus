package cron

import (
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// InitDataHistory TODO:
func InitDataHistory() {
	go func() {
		d := time.Duration(g.COLLECT_INTERVAL) * time.Second
		for range time.Tick(d) {
			funcs.UpdateCPUStats()
			funcs.UpdateDiskStats()
		}
	}()
}

// Collect TODO:
func Collect() {
	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*cmodel.MetricValue) {
	d := time.Second * time.Duration(sec)
	for range time.Tick(d) {
		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		mvs := []*cmodel.MetricValue{}
		ignoreMetrics := g.Config().IgnoreMetrics

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				if b, ok := ignoreMetrics[mv.Metric]; ok && b {
					continue
				} else {
					mvs = append(mvs, mv)
				}
			}
		}

		now := time.Now().Unix()
		for j := 0; j < len(mvs); j++ {
			mvs[j].Step = sec
			mvs[j].Endpoint = hostname
			mvs[j].Timestamp = now
		}

		g.SendToTransfer(mvs)
	}
}

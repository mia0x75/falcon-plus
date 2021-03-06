package cron

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/hbs"
)

// SyncBuiltinMetrics TODO:
func SyncBuiltinMetrics() {
	if len(g.Config().Heartbeat.Addrs) > 0 {
		go syncBuiltinMetrics()
	}
}

func syncBuiltinMetrics() {
	timestamp := int64(-1)
	checksum := "nil"

	d := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	for range time.Tick(d) {
		var ports = []int64{}
		var paths = []string{}
		var sources = make(map[string]string)
		var procs = make(map[string]map[int]string)
		var urls = make(map[string]string)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := cm.AgentHeartbeatRequest{
			Hostname: hostname,
			Checksum: checksum,
		}

		var resp cm.BuiltinMetricResponse
		err = g.HbsClient.Call("Agent.BuiltinMetrics", req, &resp)
		if err != nil {
			log.Errorf("[E] %v", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		if resp.Checksum == checksum {
			continue
		}

		timestamp = resp.Timestamp
		checksum = resp.Checksum

		for _, metric := range resp.Metrics {
			if metric.Metric == g.URL_CHECK_HEALTH {
				arr := strings.Split(metric.Tags, ",")
				if len(arr) != 2 {
					continue
				}
				url := strings.Split(arr[0], "=")
				if len(url) != 2 {
					continue
				}
				stime := strings.Split(arr[1], "=")
				if len(stime) != 2 {
					continue
				}
				if _, err := strconv.ParseInt(stime[1], 10, 64); err == nil {
					urls[url[1]] = stime[1]
				} else {
					log.Errorf("[E] metric ParseInt timeout failed: %v", err)
				}
				continue
			}

			if metric.Metric == g.NET_PORT_LISTEN {
				arr := strings.Split(metric.Tags, "=")
				if len(arr) != 2 {
					continue
				}

				if port, err := strconv.ParseInt(arr[1], 10, 64); err == nil {
					ports = append(ports, port)
				} else {
					log.Errorf("[E] metrics ParseInt failed: %v", err)
				}
				continue
			}

			if metric.Metric == g.DU_BS {
				arr := strings.Split(metric.Tags, "=")
				if len(arr) != 2 {
					continue
				}

				paths = append(paths, strings.TrimSpace(arr[1]))
				continue
			}

			if metric.Metric == g.FS_FILE_CHECKSUM {
				arr := strings.Split(metric.Tags, "=")
				if len(arr) != 2 {
					continue
				}
				if arr[0] == "source" {
					sources[arr[1]] = strings.TrimSpace(arr[1])
				} else {
					log.Errorf("[E] invalid tag: %s", arr[0])
					continue
				}
				continue
			}

			if metric.Metric == g.PROC_NUM {
				arr := strings.Split(metric.Tags, ",")

				tmpMap := make(map[int]string)

				for i := 0; i < len(arr); i++ {
					if strings.HasPrefix(arr[i], "name=") {
						tmpMap[1] = strings.TrimSpace(arr[i][5:])
					} else if strings.HasPrefix(arr[i], "cmdline=") {
						tmpMap[2] = strings.TrimSpace(arr[i][8:])
					}
				}

				procs[metric.Tags] = tmpMap
				continue
			}
		}

		if !cmp.Equal(urls, hbs.ReportUrls()) {
			hbs.CacheReportUrls(urls)
		}
		if !cmp.Equal(ports, hbs.ReportPorts()) {
			hbs.CacheReportPorts(ports)
		}
		if !cmp.Equal(procs, hbs.ReportProcs()) {
			hbs.CacheReportProcs(procs)
		}
		if !cmp.Equal(paths, hbs.ReportPaths()) {
			hbs.CacheReportPaths(paths)
		}
	}
}

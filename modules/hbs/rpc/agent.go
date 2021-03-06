package rpc

import (
	"bytes"
	"sort"
	"strings"
	"time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

// Agent TODO:
type Agent int

// MinePlugins TODO:
func (t *Agent) MinePlugins(args cm.AgentHeartbeatRequest, reply *cm.AgentPluginsResponse) error {
	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()

	return nil
}

// ReportStatus agent上报自身状态
func (t *Agent) ReportStatus(args *cm.AgentReportRequest, reply *cm.SimpleRPCResponse) error {
	if args.Hostname == "" {
		reply.Code = 1
		return nil
	}

	cache.Agents.Put(args)

	return nil
}

// TrustableIps 需要checksum一下来减少网络开销？其实白名单通常只会有一个或者没有，无需checksum
func (t *Agent) TrustableIps(args *cm.NullRPCRequest, ips *string) error {
	*ips = strings.Join(g.Config().Trustable, ",")
	return nil
}

// BuiltinMetrics agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) BuiltinMetrics(args *cm.AgentHeartbeatRequest, reply *cm.BuiltinMetricResponse) error {
	if args.Hostname == "" {
		return nil
	}

	metrics, err := cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []*cm.BuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}

// DigestBuiltinMetrics TODO:
func DigestBuiltinMetrics(items []*cm.BuiltinMetric) string {
	sort.Sort(cm.BuiltinMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return cu.Md5(buf.String())
}

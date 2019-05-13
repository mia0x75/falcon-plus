package g

import (
	"bytes"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

// LocalIP TODO:
var LocalIP string

// InitLocalIP TODO:
func InitLocalIP() {
	for _, addr := range Config().Heartbeat.Addrs {
		conn, err := net.DialTimeout("tcp", addr, time.Second*10)
		if err != nil {
			log.Errorf("[E] connect to heartbeat server %s failed", addr)
		} else {
			defer conn.Close()
			LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
			break
		}
	}
}

// TODO:
var (
	HbsClient *SingleConnRPCClient
)

// InitRPCClients TODO:
func InitRPCClients() {
	if len(Config().Heartbeat.Addrs) > 0 {
		HbsClient = &SingleConnRPCClient{
			RPCServers: Config().Heartbeat.Addrs,
			Timeout:    time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	} else {
		// TODO: panic
	}
}

// SendToTransfer TODO:
func SendToTransfer(metrics []*cmodel.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	dt := Config().DefaultTags
	if len(dt) > 0 {
		var buf bytes.Buffer
		list := []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			list = append(list, buf.String())
		}
		defaultTags := strings.Join(list, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = defaultTags
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(defaultTags)
				metrics[i].Tags = buf.String()
			}
		}
	}
	for _, m := range metrics {
		log.Debugf("[D] => Metric %v", m)
	}

	log.Debugf("[D] => <Total=%d> %v", len(metrics), metrics[0])

	var resp cmodel.TransferResponse
	SendMetrics(metrics, &resp)

	log.Debugf("[D] <= %v", &resp)
}

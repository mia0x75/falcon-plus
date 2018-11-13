package g

import (
	"bytes"
	"net"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

var LocalIp string

func InitLocalIp() {
	for _, addr := range Config().Heartbeat.Addrs {
		conn, err := net.DialTimeout("tcp", addr, time.Second*10)
		if err != nil {
			log.Errorf("[E] connect to heartbeat server %s failed", addr)
		} else {
			defer conn.Close()
			LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
			break
		}
	}
}

var (
	HbsClient *SingleConnRpcClient
)

func InitRpcClients() {
	if len(Config().Heartbeat.Addrs) > 0 {
		HbsClient = &SingleConnRpcClient{
			RpcServers: Config().Heartbeat.Addrs,
			Timeout:    time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	} else {
		// TODO: panic
	}
}

func SendToTransfer(metrics []*cmodel.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	dt := Config().DefaultTags
	if len(dt) > 0 {
		var buf bytes.Buffer
		default_tags_list := []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			default_tags_list = append(default_tags_list, buf.String())
		}
		default_tags := strings.Join(default_tags_list, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = default_tags
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(default_tags)
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

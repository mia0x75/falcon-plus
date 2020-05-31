package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	cron "github.com/toolkits/cron"
	ntime "github.com/toolkits/time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

var collectorCron = cron.New()

// Start 启动服务
func Start() {
	if !g.Config().Collector.Enabled {
		log.Info("[I] collector.Start warning, not enable")
		return
	}

	if g.Config().Collector.Agent == "" {
		return
	}
	if g.Config().Collector.Pattern == "" {
		return
	}
	// Start
	go startCollectorCron()
	log.Info("[I] collector.Start, ok")
}

func startCollectorCron() {
	collectorCron.AddFuncCC("0 * * * * ?", func() { collect() }, 1)
	collectorCron.Start()
}

func collect() {
	startTs := time.Now().Unix()
	_collect()
	endTs := time.Now().Unix()
	log.Infof("[I] collect, start %s, ts %ds\n", ntime.FormatTs(startTs), endTs-startTs)

	// Statistics
	proc.CollectorCronCnt.Incr()
}

func _collect() {
	for _, host := range g.Config().Collector.Cluster {
		ts := time.Now().Unix()
		jsonList := make([]*cm.JSONMetaData, 0)

		// Get statistics via http-get
		hostInfo := strings.Split(host, ",") // "module,hostname:port"
		if len(hostInfo) != 2 {
			continue
		}
		hostModule := hostInfo[0]
		hostNamePort := hostInfo[1]

		hostNamePortList := strings.Split(hostNamePort, ":")
		if len(hostNamePortList) != 2 {
			continue
		}
		hostName := hostNamePortList[0]
		hostPort := hostNamePortList[1]

		tags := "port=" + hostPort
		srcURL := fmt.Sprintf(g.Config().Collector.Pattern, hostNamePort)
		client := cu.NewHttp(srcURL)
		client.SetUserAgent("collector.get")
		headers := map[string]string{
			"Connection": "close",
		}
		client.SetHeaders(headers)
		body, err := client.Get()
		if err != nil {
			log.Infof("[I] %s, get statistics error: %v", hostNamePort, err)
			continue
		}

		var data Dto
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Infof("[I] %s, get statistics error: %v", hostNamePort, err)
			continue
		}

		for _, item := range data.Data {
			if item["Name"] == nil {
				continue
			}
			itemName := item["Name"].(string)
			if item["Cnt"] != nil {
				var jmdCnt cm.JSONMetaData
				jmdCnt.Endpoint = hostName
				jmdCnt.Metric = fmt.Sprintf("%s.stats.%s", hostModule, itemName)
				jmdCnt.Timestamp = ts
				jmdCnt.Step = 60
				jmdCnt.Value = int64(item["Cnt"].(float64))
				jmdCnt.CounterType = "GAUGE"
				jmdCnt.Tags = tags
				jsonList = append(jsonList, &jmdCnt)
			}

			if item["Qps"] != nil {
				var jmdQPS cm.JSONMetaData
				jmdQPS.Endpoint = hostName
				jmdQPS.Metric = fmt.Sprintf("%s.stats.%s.Qps", hostModule, itemName)
				jmdQPS.Timestamp = ts
				jmdQPS.Step = 60
				jmdQPS.Value = int64(item["Qps"].(float64))
				jmdQPS.CounterType = "GAUGE"
				jmdQPS.Tags = tags
				jsonList = append(jsonList, &jmdQPS)
			}
		}

		// Format result
		err = sendToTransfer(jsonList, g.Config().Collector.Agent)
		if err != nil {
			log.Infof("[I] %s send to transfer error: %v", hostNamePort, err)
		}
	}
}

func sendToTransfer(items []*cm.JSONMetaData, destURL string) error {
	if len(items) < 1 {
		return nil
	}

	// Format result
	jsonBody, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("json.Marshal failed with %v", err)
	}

	// Send via http-post
	client := cu.NewHttp(destURL)
	client.SetUserAgent("collector.post")
	headers := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Connection":   "close",
	}
	client.SetHeaders(headers)
	if _, err = client.Post(jsonBody); err != nil {
		return fmt.Errorf("post to %s, resquest failed with %v", destURL, err)
	}

	return nil
}

// Dto TODO:
type Dto struct {
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}

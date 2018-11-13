package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/common/model"
	cron "github.com/toolkits/cron"
	ntime "github.com/toolkits/time"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

var collectorCron = cron.New()

func Start() {
	if !g.Config().Collector.Enabled {
		log.Info("[I] collector.Start warning, not enable")
		return
	}

	// init url
	if g.Config().Collector.Agent == "" {
		return
	}
	if g.Config().Collector.Pattern == "" {
		return
	}
	// start
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

	// statistics
	proc.CollectorCronCnt.Incr()
}

func _collect() {
	for _, host := range g.Config().Collector.Cluster {
		ts := time.Now().Unix()
		jsonList := make([]*cmodel.JsonMetaData, 0)

		// get statistics by http-get
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
		srcUrl := fmt.Sprintf(g.Config().Collector.Pattern, hostNamePort)
		client := cutils.NewHttp(srcUrl)
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
				var jmdCnt cmodel.JsonMetaData
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
				var jmdQps cmodel.JsonMetaData
				jmdQps.Endpoint = hostName
				jmdQps.Metric = fmt.Sprintf("%s.stats.%s.Qps", hostModule, itemName)
				jmdQps.Timestamp = ts
				jmdQps.Step = 60
				jmdQps.Value = int64(item["Qps"].(float64))
				jmdQps.CounterType = "GAUGE"
				jmdQps.Tags = tags
				jsonList = append(jsonList, &jmdQps)
			}
		}

		// format result
		err = sendToTransfer(jsonList, g.Config().Collector.Agent)
		if err != nil {
			log.Infof("[I] %s send to transfer error: %v", hostNamePort, err)
		}
	}
}

func sendToTransfer(items []*cmodel.JsonMetaData, destUrl string) error {
	if len(items) < 1 {
		return nil
	}

	// format result
	jsonBody, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("json.Marshal failed with %v", err)
	}

	// send by http-post
	client := cutils.NewHttp(destUrl)
	client.SetUserAgent("collector.post")
	headers := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Connection":   "close",
	}
	client.SetHeaders(headers)
	if _, err = client.Post(jsonBody); err != nil {
		return fmt.Errorf("post to %s, resquest failed with %v", destUrl, err)
	}

	return nil
}

type Dto struct {
	Msg  string                   `json:"msg"`
	Data []map[string]interface{} `json:"data"`
}

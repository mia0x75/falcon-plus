package index

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	cron "github.com/toolkits/cron"
	ntime "github.com/toolkits/time"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/exporter/g"
	"github.com/open-falcon/falcon-plus/modules/exporter/proc"
)

const (
	destUrlFmt = "http://%s/index/updateAll"
)

var (
	indexUpdateAllCron = cron.New()
)

// 启动 索引全量更新 定时任务
func StartIndexUpdateAllTask() {
	for graphAddr, cronSpec := range g.Config().Index.Cluster {
		ga := graphAddr
		indexUpdateAllCron.AddFuncCC(cronSpec, func() { UpdateIndexOfOneGraph(ga, "cron") }, 1)
	}

	indexUpdateAllCron.Start()
}

// 手动触发全量更新
func UpdateAllIndex() {
	for graphAddr := range g.Config().Index.Cluster {
		UpdateIndexOfOneGraph(graphAddr, "manual")
	}
}

func UpdateIndexOfOneGraph(graphAddr string, src string) {
	startTs := time.Now().Unix()
	err := updateIndexOfOneGraph(graphAddr)
	endTs := time.Now().Unix()

	// statistics
	proc.IndexUpdateCnt.Incr()
	if err == nil {
		log.Infof(
			"[I] index update ok, %s, %s, start %s, ts %ds",
			src,
			graphAddr,
			ntime.FormatTs(startTs),
			endTs-startTs,
		)
	} else {
		proc.IndexUpdateErrorCnt.Incr()
		log.Errorf(
			"[E] index update error, %s, %s, start %s, ts %ds, reason %v",
			src,
			graphAddr,
			ntime.FormatTs(startTs),
			endTs-startTs,
			err,
		)
	}
}

func updateIndexOfOneGraph(hostNamePort string) error {
	if hostNamePort == "" {
		return fmt.Errorf("index update error, bad host")
	}

	destUrl := fmt.Sprintf(destUrlFmt, hostNamePort)

	client := cu.NewHttp(destUrl)
	client.SetUserAgent(fmt.Sprintf("index.update.%s", hostNamePort))
	headers := map[string]string{
		"Connection": "close",
	}
	client.SetHeaders(headers)
	body, err := client.Get()
	if err != nil {
		log.Errorf("[E] %s, index update error: %v", hostNamePort, err)
		return err
	}

	var data Dto
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Errorf("[E] %s, index update error: %v", hostNamePort, err)
		return err
	}

	if data.Data != "ok" {
		log.Infof("[I] %s, index update error, bad result: %s", hostNamePort, data.Data)
		return err
	}

	return nil
}

type Dto struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

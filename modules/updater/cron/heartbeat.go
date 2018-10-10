package cron

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/updater/g"
	"github.com/open-falcon/ops-common/model"
	"github.com/open-falcon/ops-common/utils"
	"github.com/toolkits/net/httplib"
)

func Heartbeat() {
	go func() {
		SleepRandomDuration()
		d := time.Duration(g.Config().Interval) * time.Second
		for range time.Tick(d) {
			heartbeat()
		}
	}()
}

func heartbeat() {
	agentDirs, err := ListAgentDirs()
	if err != nil {
		return
	}

	hostname, err := utils.Hostname(g.Config().Hostname)
	if err != nil {
		return
	}

	heartbeatRequest := BuildHeartbeatRequest(hostname, agentDirs)
	log.Debugln("====>>>>")
	log.Debugln(heartbeatRequest)

	bs, err := json.Marshal(heartbeatRequest)
	if err != nil {
		log.Println("encode heartbeat request fail", err)
		return
	}

	url := fmt.Sprintf("http://%s/hosts", g.Config().Server)

	httpRequest := httplib.Post(url).SetTimeout(time.Second*10, time.Minute)
	httpRequest.Body(bs)
	httpResponse, err := httpRequest.Bytes()
	if err != nil {
		log.Printf("curl %s fail %v", url, err)
		return
	}

	var heartbeatResponse model.HeartbeatResponse
	err = json.Unmarshal(httpResponse, &heartbeatResponse)
	if err != nil {
		log.Println("decode heartbeat response fail", err)
		return
	}

	log.Debugln("<<<<====")
	log.Debugln(heartbeatResponse)

	HandleHeartbeatResponse(&heartbeatResponse)

}

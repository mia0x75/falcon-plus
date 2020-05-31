package cron

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

// ConsumeIM 处理队列
func ConsumeIM() {
	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			L := redi.PopAllIM()
			if len(L) == 0 {
				continue
			}
			SendIMList(L)
		}
	}()
}

// SendIMList 处理IM告警队列
func SendIMList(L []*g.AlarmDto) {
	for _, im := range L {
		IMWorkerChan <- 1
		go SendIM(im)
	}
}

// SendIM 发送IM告警
func SendIM(im *g.AlarmDto) {
	defer func() {
		<-IMWorkerChan
	}()

	url := g.Config().API.IM
	if strings.TrimSpace(url) != "" {
		if data, err := json.Marshal(im); err != nil {
			log.Errorf("[E] %v", err)
		} else {
			resp, err := cu.Post(url, data)
			if err != nil {
				log.Errorf("[E] send im fail, content: %v, error: %v", im, err)
			}
			log.Debugf("[D] send im: %v, resp: %v, url: %s", im, resp, url)
		}
	} else {
		log.Debugf("[D] im url: %s is blank, SKIP", url)
	}
}

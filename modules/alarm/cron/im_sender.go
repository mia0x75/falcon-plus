package cron

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

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

func SendIMList(L []*g.AlarmDto) {
	for _, im := range L {
		IMWorkerChan <- 1
		go SendIM(im)
	}
}

func SendIM(im *g.AlarmDto) {
	defer func() {
		<-IMWorkerChan
	}()

	url := g.Config().Api.IM
	if strings.TrimSpace(url) != "" {
		if data, err := json.Marshal(im); err != nil {
			log.Errorf("[ERROF] %v", err)
			return
		} else {
			resp, err := cutils.Post(url, data)
			if err != nil {
				log.Errorf("[E] send im fail, content: %v, error: %v", im, err)
			}
			log.Debugf("[D] send im: %v, resp: %v, url: %s", im, resp, url)
		}
	} else {
		log.Debugf("[D] im url: %s is blank, SKIP", url)
	}
}

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

func ConsumeSms() {
	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			L := redi.PopAllSms()
			if len(L) == 0 {
				continue
			}
			SendSmsList(L)
		}
	}()
}

func SendSmsList(L []*g.AlarmDto) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

func SendSms(sms *g.AlarmDto) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config().Api.Sms
	if strings.TrimSpace(url) != "" {
		if data, err := json.Marshal(sms); err != nil {
			log.Error(err)
			return
		} else {
			resp, err := cutils.Post(url, data)
			if err != nil {
				log.Errorf("send sms fail, content:%v, error:%v", sms, err)
			}
			log.Debugf("send sms:%v, resp:%v, url:%s", sms, resp, url)
		}
	} else {
		log.Debugf("sms url:%s is blank, SKIP", url)
	}
}

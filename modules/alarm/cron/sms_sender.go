package cron

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

// ConsumeSms 处理队列
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

// SendSmsList 处理短信告警队列
func SendSmsList(L []*g.AlarmDto) {
	for _, sms := range L {
		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

// SendSms 发送短信告警
func SendSms(sms *g.AlarmDto) {
	defer func() {
		<-SmsWorkerChan
	}()

	addr := g.Config().API.Sms
	for {
		// Blank
		if strings.TrimSpace(addr) == "" {
			log.Debugf("[D] sms url: %s is blank, SKIP", addr)
			break
		}
		// URL validation
		if _, err := url.ParseRequestURI(addr); err != nil {
			log.Errorf("[E] %s is not a valid url.", addr)
			break
		}
		if data, err := json.Marshal(sms); err != nil {
			log.Errorf("[E] %v", err)
			break
		} else {
			resp, err := cutils.Post(addr, data)
			if err != nil {
				log.Errorf("[E] send sms fail, content: %v, error: %v", sms, err)
				break
			}
			log.Debugf("[D] send sms: %v, resp: %v, url: %s", sms, resp, addr)
		}
		break
	}
}

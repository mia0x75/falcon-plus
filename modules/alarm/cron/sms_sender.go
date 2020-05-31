package cron

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

// ConsumeSMS 处理队列
func ConsumeSMS() {
	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			L := redi.PopAllSMS()
			if len(L) == 0 {
				continue
			}
			SendSMSList(L)
		}
	}()
}

// SendSMSList 处理短信告警队列
func SendSMSList(L []*g.AlarmDto) {
	for _, s := range L {
		SMSWorkerChan <- 1
		go SendSMS(s)
	}
}

// SendSMS 发送短信告警
func SendSMS(s *g.AlarmDto) {
	defer func() {
		<-SMSWorkerChan
	}()

	addr := g.Config().API.SMS
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
		if data, err := json.Marshal(s); err != nil {
			log.Errorf("[E] %v", err)
			break
		} else {
			resp, err := cu.Post(addr, data)
			if err != nil {
				log.Errorf("[E] send sms fail, content: %v, error: %v", s, err)
				break
			}
			log.Debugf("[D] send sms: %v, resp: %v, url: %s", s, resp, addr)
		}
		break
	}
}

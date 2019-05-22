package cron

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

// ConsumeMail 处理队列
func ConsumeMail() {
	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			L := redi.PopAllMail()
			if len(L) == 0 {
				continue
			}
			SendMailList(L)
		}
	}()
}

// SendMailList 处理邮件告警队列
func SendMailList(L []*g.AlarmDto) {
	for _, mail := range L {
		MailWorkerChan <- 1
		go SendMail(mail)
	}
}

// SendMail 发送邮件告警
func SendMail(mail *g.AlarmDto) {
	defer func() {
		<-MailWorkerChan
	}()

	url := g.Config().API.Mail
	if strings.TrimSpace(url) != "" {
		if data, err := json.Marshal(mail); err != nil {
			log.Errorf("[E] %v", err)
		} else {
			resp, err := cutils.Post(url, data)
			if err != nil {
				log.Errorf("[E] send mail fail, content: %v, error: %v", mail, err)
			}
			log.Debugf("[D] send mail: %v, resp: %v, url: %s", mail, resp, url)
		}
	} else {
		log.Debugf("[D] mail url: %s is blank, SKIP", url)
	}
}

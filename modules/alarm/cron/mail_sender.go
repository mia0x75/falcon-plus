package cron

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
)

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

func SendMailList(L []*model.Mail) {
	for _, mail := range L {
		MailWorkerChan <- 1
		go SendMail(mail)
	}
}

func SendMail(mail *model.Mail) {
	defer func() {
		<-MailWorkerChan
	}()

	url := g.Config().Api.Mail
	if strings.TrimSpace(url) != "" {
		r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
		r.Param("tos", mail.Tos)
		r.Param("subject", mail.Subject)
		r.Param("content", mail.Content)
		resp, err := r.String()
		if err != nil {
			log.Errorf("send mail fail, receiver:%s, subject:%s, cotent:%s, error:%v", mail.Tos, mail.Subject, mail.Content, err)
		}
		log.Debugf("send mail:%v, resp:%v, url:%s", mail, resp, url)
	} else {
		log.Debugf("mail url:%s is blank, SKIP", url)
	}
}

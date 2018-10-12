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

func SendIMList(L []*model.IM) {
	for _, im := range L {
		IMWorkerChan <- 1
		go SendIM(im)
	}
}

func SendIM(im *model.IM) {
	defer func() {
		<-IMWorkerChan
	}()

	url := g.Config().Api.IM
	log.Debugf("send im via %s", url)
	if strings.TrimSpace(url) != "" {
		r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
		r.Param("tos", im.Tos)
		r.Param("content", im.Content)
		resp, err := r.String()
		if err != nil {
			log.Errorf("send im fail, tos:%s, cotent:%s, error:%v", im.Tos, im.Content, err)
		}
		log.Debugf("send im:%v, resp:%v, url:%s", im, resp, url)
	} else {
		log.Debugf("im url:%s is blank, SKIP", url)
	}
}

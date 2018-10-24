package redi

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func lpush(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Error("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteSms(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write sms to queue, sms:%v, queue:%s", content, g.SMS_QUEUE_NAME)
	lpush(g.SMS_QUEUE_NAME, string(bs))
}

func WriteIM(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write im to queue, im:%v, queue:%s", content, g.IM_QUEUE_NAME)
	lpush(g.IM_QUEUE_NAME, string(bs))
}

func WriteMail(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write mail to queue, mail:%v, queue:%s", content, g.MAIL_QUEUE_NAME)
	lpush(g.MAIL_QUEUE_NAME, string(bs))
}

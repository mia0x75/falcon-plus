package redi

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func lpush(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Errorf("[E] LPUSH redis %s fail: %v dto: %s", queue, err, message)
	}
}

// WriteSMS 短信告警排队
func WriteSMS(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	log.Debugf("[D] write sms to queue, sms: %v, queue: %s", content, g.Config().Queue.InstantQueues.SMSQueue)
	lpush(g.Config().Queue.InstantQueues.SMSQueue, string(bs))
}

// WriteIM 即时消息告警排队
func WriteIM(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	log.Debugf("[D] write im to queue, im: %v, queue: %s", content, g.Config().Queue.InstantQueues.IMQueue)
	lpush(g.Config().Queue.InstantQueues.IMQueue, string(bs))
}

// WriteMail 电子邮件告警排队
func WriteMail(content *g.AlarmDto) {
	if content.Subscriber == nil {
		return
	}
	if len(content.Subscriber) == 0 {
		return
	}

	bs, err := json.Marshal(content)
	if err != nil {
		log.Errorf("[E] %v", err)
		return
	}

	log.Debugf("[D] write mail to queue, mail: %v, queue: %s", content, g.Config().Queue.InstantQueues.MailQueue)
	lpush(g.Config().Queue.InstantQueues.MailQueue, string(bs))
}

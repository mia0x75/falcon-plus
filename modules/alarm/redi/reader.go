package redi

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

// PopAllSMS 短信告警读取
func PopAllSMS() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.InstantQueues.SMSQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Errorf("[E] %v", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var sms *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &sms)
		if err != nil {
			log.Errorf("[E] reply: %s, error: %v", reply, err)
			continue
		}

		ret = append(ret, sms)
	}

	return ret
}

// PopAllIM 即时消息告警读取
func PopAllIM() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.InstantQueues.IMQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Errorf("[E] %v", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var im *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &im)
		if err != nil {
			log.Errorf("[E] reply: %s, error: %v", reply, err)
			continue
		}

		ret = append(ret, im)
	}

	return ret
}

// PopAllMail 邮件告警读取
func PopAllMail() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.InstantQueues.MailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Errorf("[E] %v", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mail *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &mail)
		if err != nil {
			log.Errorf("[E] reply: %s, error: %v", reply, err)
			continue
		}
		ret = append(ret, mail)
	}

	return ret
}

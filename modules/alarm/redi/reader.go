package redi

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func PopAllSms() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.SMS_QUEUE_NAME

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

func PopAllIM() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.IM_QUEUE_NAME

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

func PopAllMail() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.MAIL_QUEUE_NAME

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

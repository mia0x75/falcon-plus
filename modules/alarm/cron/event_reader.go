package cron

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func ReadHighEvent() {
	queues := g.Config().Queue.HighQueues
	if len(queues) == 0 {
		return
	}

	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			event, err := popEvent(queues)
			if err != nil {
				continue
			}
			consume(event, true)
		}
	}()
}

func ReadLowEvent() {
	queues := g.Config().Queue.LowQueues
	if len(queues) == 0 {
		return
	}

	go func() {
		d := time.Duration(1) * time.Second
		for range time.Tick(d) {
			event, err := popEvent(queues)
			if err != nil {
				continue
			}
			consume(event, false)
		}
	}()
}

func popEvent(queues []string) (*cmodel.Event, error) {
	count := len(queues)
	params := make([]interface{}, count+1)
	for i := 0; i < count; i++ {
		params[i] = queues[i]
	}
	// set timeout 0
	params[count] = 0

	rc := g.RedisConnPool.Get()
	if rc == nil {
		log.Warn("[W] cannot get redis connection")
		return nil, errors.New("get redis connection failed")
	}
	defer rc.Close()

	reply, err := redis.Strings(rc.Do("BRPOP", params...))
	if err != nil {
		log.Errorf("[E] get alarm event from redis fail: %v", err)
		return nil, err
	}

	var event cmodel.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Errorf("[E] parse alarm event fail: %v", err)
		return nil, err
	}

	log.Debugf("[D] pop event: %s", event.String())

	//insert event into database
	eventmodel.InsertEvent(&event)
	// events no longer saved in memory

	return &event, nil
}

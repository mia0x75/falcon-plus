package cron

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func ReadHighEvent() {
	queues := g.Config().Redis.HighQueues
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
	queues := g.Config().Redis.LowQueues
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
	defer rc.Close()
	if rc == nil {
		return nil, errors.New("get redis connection failed")
	}

	reply, err := redis.Strings(rc.Do("BRPOP", params...))
	if err != nil {
		log.Errorf("get alarm event from redis fail: %v", err)
		return nil, err
	}

	var event cmodel.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Errorf("parse alarm event fail: %v", err)
		return nil, err
	}

	log.Debugf("pop event: %s", event.String())

	//insert event into database
	eventmodel.InsertEvent(&event)
	// events no longer saved in memory

	return &event, nil
}

package cron

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/api/app/model"
)

// ReadHighEvent 读取高优先级的告警队列
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

// ReadLowEvent 读取低优先级的告警队列
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

func popEvent(queues []string) (*cm.Event, error) {
	count := len(queues)
	params := make([]interface{}, count+1)
	for i := 0; i < count; i++ {
		params[i] = queues[i]
	}
	// Set timeout 0
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

	var event cm.Event
	err = json.Unmarshal([]byte(reply[1]), &event)
	if err != nil {
		log.Errorf("[E] parse alarm event fail: %v", err)
		return nil, err
	}

	log.Debugf("[D] pop event: %s", event.String())

	InsertEvent(&event)

	return &event, nil
}

func insertEvent(db *gorm.DB, eve *cm.Event) {
	var status int
	if status = 0; eve.Status == "OK" {
		status = 1
	}
	event := model.Event{
		CaseID: eve.ID,
		Step:   eve.CurrentStep,
		Cond:   fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		Status: status,
	}
	if err := db.Create(&event).Error; err != nil {
		log.Errorf("[E] insert event to db fail, error: %v", err)
	} else {
		log.Debugf("[D] insert event to db succ, last_insert_id: %d", event.ID)
	}
	return
}

// InsertEvent TODO:
func InsertEvent(eve *cm.Event) {
	db := g.Con()
	dt := db
	var cases []model.Case
	db.Where("id = ?", eve.ID).Find(&cases)
	log.Debugf("[D] events: %v", eve)
	log.Debugf("[D] expression is null: %v", eve.Expression == nil)
	if len(cases) == 0 {
		c := model.Case{
			ID:           eve.ID,
			Endpoint:     eve.Endpoint,
			Metric:       counterGen(eve.Metric(), cu.SortedTags(eve.PushedTags)),
			Func:         eve.Func(),
			Cond:         fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			Note:         eve.Note(),
			MaxStep:      eve.MaxStep(),
			CurrentStep:  eve.CurrentStep,
			Priority:     eve.Priority(),
			Status:       eve.Status,
			ExpressionID: int64(eve.ExpressionID()),
			StrategyID:   int64(eve.StrategyID()),
			TemplateID:   int64(eve.TemplateID()),
		}
		// Create cases
		db.Create(&c)
	} else {
		c := model.Case{
			ID:           eve.ID,
			Func:         eve.Func(),
			Cond:         fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			Note:         eve.Note(),
			MaxStep:      eve.MaxStep(),
			CurrentStep:  eve.CurrentStep,
			Priority:     eve.Priority(),
			Status:       eve.Status,
			ExpressionID: int64(eve.ExpressionID()),
			StrategyID:   int64(eve.StrategyID()),
			TemplateID:   int64(eve.TemplateID()),
		}

		// Re-open the case
		if cases[0].ProcessStatus == "resolved" || cases[0].ProcessStatus == "ignored" {
			c.ProcessNote = 0
			c.ProcessStatus = "unresolved"
		}

		if eve.CurrentStep == 1 {
			// Update start time of cases
			c.CreateAt = time.Now().Unix()
		}
		db.Where("id = ?", eve.ID).Update(&c)
	}
	log.Debugf("[D] %v", dt.Error)
	// Insert case
	insertEvent(db, eve)
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}

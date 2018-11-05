package cron

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

func CombineSms() {
	go func() {
		d := time.Duration(1) * time.Minute
		// 每分钟读取处理一次
		for range time.Tick(d) {
			combineSms()
		}
	}()
}

func CombineMail() {
	go func() {
		d := time.Duration(1) * time.Minute
		// 每分钟读取处理一次
		for range time.Tick(d) {
			combineMail()
		}
	}()
}

func CombineIM() {
	go func() {
		d := time.Duration(1) * time.Minute
		// 每分钟读取处理一次
		for range time.Tick(d) {
			combineIM()
		}
	}()
}

func combineMail() {
	dtos := popAllMailDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*g.AlarmDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Subscriber, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*g.AlarmDto{dtos[i]}
		}
	}

	// 不要在这处理，继续写回redis，否则重启alarm很容易丢数据
	for _, arr := range dtoMap {
		if arr[0].Subscriber == nil {
			continue
		}
		if len(arr[0].Subscriber) == 0 {
			continue
		}
		size := len(arr)
		arr[0].Occur = size
		redi.WriteMail(arr[0])
	}
}

func combineIM() {
	dtos := popAllImDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*g.AlarmDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Subscriber, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*g.AlarmDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		if arr[0].Subscriber == nil {
			continue
		}
		if len(arr[0].Subscriber) == 0 {
			continue
		}
		size := len(arr)
		arr[0].Occur = size
		redi.WriteIM(arr[0])
	}
}

func combineSms() {
	dtos := popAllSmsDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*g.AlarmDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Subscriber, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*g.AlarmDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		if arr[0].Subscriber == nil {
			continue
		}
		if len(arr[0].Subscriber) == 0 {
			continue
		}
		size := len(arr)
		arr[0].Occur = size
		redi.WriteSms(arr[0])
	}
}

func popAllSmsDto() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get SmsDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var smsDto *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &smsDto)
		if err != nil {
			log.Errorf("json unmarshal SmsDto: %s fail: %v", reply, err)
			continue
		}
		ret = append(ret, smsDto)
	}

	return ret
}

func popAllMailDto() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get MailDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mailDto *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &mailDto)
		if err != nil {
			log.Errorf("json unmarshal MailDto: %s fail: %v", reply, err)
			continue
		}
		ret = append(ret, mailDto)
	}

	return ret
}

func popAllImDto() []*g.AlarmDto {
	var ret []*g.AlarmDto
	queue := g.Config().Queue.UserIMQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get ImDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var imDto *g.AlarmDto
		err = json.Unmarshal([]byte(reply), &imDto)
		if err != nil {
			log.Errorf("json unmarshal imDto: %s fail: %v", reply, err)
			continue
		}
		ret = append(ret, imDto)
	}

	return ret
}

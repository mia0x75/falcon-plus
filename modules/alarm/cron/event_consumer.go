package cron

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

func consume(event *cmodel.Event, isHigh bool) {
	actionID := event.ActionId()
	if actionID <= 0 {
		return
	}

	action := api.GetAction(actionID)
	if action == nil {
		return
	}

	if isHigh {
		consumeHighEvents(event, action)
	} else {
		consumeLowEvents(event, action)
	}
}

// 高优先级的不做报警合并
func consumeHighEvents(event *cmodel.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	phones, mails, ims := api.ParseTeams(action.Uic)
	// <=P2 才发送短信
	if event.Priority() < 3 {
		if len(phones) > 0 {
			content := GenerateSmsContent(phones, event)
			content.Uic = action.Uic
			redi.WriteSms(content)
		}
	}

	if len(ims) > 0 {
		content := GenerateIMContent(ims, event)
		content.Uic = action.Uic
		redi.WriteIM(content)
	}
	if len(mails) > 0 {
		content := GenerateMailContent(mails, event)
		content.Uic = action.Uic
		redi.WriteMail(content)
	}
}

// 低优先级的做报警合并
func consumeLowEvents(event *cmodel.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	// <=P2 才发送短信
	if event.Priority() < 3 {
		ParseUserSms(event, action)
	}

	ParseUserIm(event, action)
	ParseUserMail(event, action)
}

// ParseUserSms TODO:
func ParseUserSms(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.LatentQueues.SmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		if user.Phone == "" {
			continue
		}
		content := GenerateSmsContent([]string{user.Phone}, event)
		content.Uic = action.Uic
		bs, err := json.Marshal(content)
		if err != nil {
			log.Errorf("[E] json marshal SmsDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("[E] LPUSH redis %s fail: %v dto: %s", queue, err, string(bs))
		}
	}
}

// ParseUserMail TODO:
func ParseUserMail(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.LatentQueues.MailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		if user.Email == "" {
			continue
		}
		content := GenerateMailContent([]string{user.Email}, event)
		content.Uic = action.Uic
		bs, err := json.Marshal(content)
		if err != nil {
			log.Errorf("[E] json marshal MailDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("[E] LPUSH redis %s fail: %v dto: %s", queue, err, string(bs))
		}
	}
}

// ParseUserIm TODO:
func ParseUserIm(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.LatentQueues.IMQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		if user.IM == "" {
			continue
		}
		content := GenerateIMContent([]string{user.IM}, event)
		content.Uic = action.Uic
		bs, err := json.Marshal(content)
		if err != nil {
			log.Errorf("[E] json marshal ImDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("[E] LPUSH redis %s fail: %v dto: %s", queue, err, string(bs))
		}
	}
}

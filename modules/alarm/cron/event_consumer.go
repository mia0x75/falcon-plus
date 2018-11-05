package cron

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
)

func consume(event *cmodel.Event, isHigh bool) {
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	action := api.GetAction(actionId)
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

func ParseUserSms(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.UserSmsQueue

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
			log.Error("json marshal SmsDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Error("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserMail(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.UserMailQueue

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
			log.Error("json marshal MailDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Error("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserIm(event *cmodel.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	queue := g.Config().Queue.UserIMQueue

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
			log.Error("json marshal ImDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Error("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

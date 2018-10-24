package cron

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func BuildContent(tos []string, event *cmodel.Event) *g.AlarmDto {
	link := g.Link(event)
	data := &g.AlarmDto{
		Status:     event.Status,
		Priority:   event.Priority(),
		Endpoint:   event.Endpoint,
		Metric:     event.Metric(),
		Tags:       cutils.SortedTags(event.PushedTags),
		Func:       event.Func(),
		LeftValue:  cutils.ReadableFloat(event.LeftValue),
		Operator:   event.Operator(),
		RightValue: cutils.ReadableFloat(event.RightValue()),
		Note:       event.Note(),
		Max:        event.MaxStep(),
		Current:    event.CurrentStep,
		Timestamp:  event.FormattedTime(),
		Link:       link,
		Subscriber: tos,
		Occur:      1,
	}
	return data
}

func GenerateSmsContent(phones []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(phones, event)
}

func GenerateMailContent(mails []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(mails, event)
}

func GenerateIMContent(ims []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(ims, event)
}

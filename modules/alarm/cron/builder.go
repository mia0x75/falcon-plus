package cron

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

// BuildContent 构建告警内容对象
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

// GenerateSmsContent 生成短信告警内容
func GenerateSmsContent(phones []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(phones, event)
}

// GenerateMailContent 生成邮件告警内容
func GenerateMailContent(mails []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(mails, event)
}

// GenerateIMContent 生成IM告警内容
func GenerateIMContent(ims []string, event *cmodel.Event) *g.AlarmDto {
	return BuildContent(ims, event)
}

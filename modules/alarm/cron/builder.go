package cron

import (
	"fmt"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func BuildCommonSMSContent(event *cmodel.Event) string {
	return fmt.Sprintf(
		"[%s][%s][][%s %s %s %s%s%s]",
		event.Status,
		event.Endpoint,
		event.Func(),
		event.Metric(),
		cutils.SortedTags(event.PushedTags),
		cutils.ReadableFloat(event.LeftValue),
		event.Operator(),
		cutils.ReadableFloat(event.RightValue()),
	)
}

func BuildCommonIMContent(event *cmodel.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		cutils.SortedTags(event.PushedTags),
		cutils.ReadableFloat(event.LeftValue),
		event.Operator(),
		cutils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *cmodel.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\nP%d\r\nEndpoint:%s\r\nMetric:%s\r\nTags:%s\r\n%s: %s%s%s\r\nNote:%s\r\nMax:%d, Current:%d\r\nTimestamp:%s\r\n%s\r\n",
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		cutils.SortedTags(event.PushedTags),
		event.Func(),
		cutils.ReadableFloat(event.LeftValue),
		event.Operator(),
		cutils.ReadableFloat(event.RightValue()),
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
	)
}

func GenerateSmsContent(event *cmodel.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *cmodel.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateIMContent(event *cmodel.Event) string {
	return BuildCommonIMContent(event)
}

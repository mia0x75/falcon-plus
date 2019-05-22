package g

import (
	"fmt"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func Link(event *cmodel.Event) string {
	tplId := event.TplId()
	if tplId != 0 {
		return fmt.Sprintf("%s/portal/template/view/%d", Config().API.Dashboard, tplId)
	}

	eid := event.ExpressionId()
	if eid != 0 {
		return fmt.Sprintf("%s/portal/expression/view/%d", Config().API.Dashboard, eid)
	}

	return ""
}

type AlarmDto struct {
	Status     string
	Priority   int
	Endpoint   string
	Metric     string
	Tags       string
	Func       string
	LeftValue  string
	Operator   string
	RightValue string
	Note       string
	Max        int
	Current    int
	Timestamp  string
	Link       string
	Occur      int
	Subscriber []string
	Uic        string
}

func (this *AlarmDto) String() string {
	return fmt.Sprintf(
		"<Content: %s, Priority:P%d, Status: %s, Value: %s, Operator: %s Threshold: %s, Occur: %d, Uic: %s, Tos: %s>",
		this.Note,
		this.Priority,
		this.Status,
		this.LeftValue,
		this.Operator,
		this.RightValue,
		this.Occur,
		this.Uic,
		this.Subscriber,
	)
}

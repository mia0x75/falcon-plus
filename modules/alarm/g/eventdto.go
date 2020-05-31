package g

import (
	"fmt"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

// Link TODO:
func Link(event *cm.Event) string {
	templateID := event.TemplateID()
	if templateID != 0 {
		return fmt.Sprintf("%s/portal/template/view/%d", Config().API.Dashboard, templateID)
	}

	eid := event.ExpressionID()
	if eid != 0 {
		return fmt.Sprintf("%s/portal/expression/view/%d", Config().API.Dashboard, eid)
	}

	return ""
}

// AlarmDto 告警数据结构
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

func (s *AlarmDto) String() string {
	return fmt.Sprintf(
		"<Content: %s, Priority:P%d, Status: %s, Value: %s, Operator: %s Threshold: %s, Occur: %d, Uic: %s, Tos: %s>",
		s.Note,
		s.Priority,
		s.Status,
		s.LeftValue,
		s.Operator,
		s.RightValue,
		s.Occur,
		s.Uic,
		s.Subscriber,
	)
}

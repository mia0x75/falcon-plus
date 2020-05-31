package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

// 机器监控和实例监控都会产生Event，共用这么一个struct
type Event struct {
	ID          string            `json:"ID"`
	Strategy    *Strategy         `json:"strategy"`
	Expression  *Expression       `json:"expression"`
	Status      string            `json:"status"` // OK or PROBLEM
	Endpoint    string            `json:"endpoint"`
	LeftValue   float64           `json:"leftValue"`
	CurrentStep int               `json:"currentStep"`
	EventTime   int64             `json:"eventTime"`
	PushedTags  map[string]string `json:"pushedTags"`
}

func (m *Event) FormattedTime() string {
	return utils.UnixTsFormat(m.EventTime)
}

func (m *Event) String() string {
	return fmt.Sprintf(
		"<Endpoint: %s, Status: %s, Strategy: %v, Expression: %v, LeftValue: %s, CurrentStep: %d, PushedTags: %v, TS: %s>",
		m.Endpoint,
		m.Status,
		m.Strategy,
		m.Expression,
		utils.ReadableFloat(m.LeftValue),
		m.CurrentStep,
		m.PushedTags,
		m.FormattedTime(),
	)
}

func (m *Event) ExpressionID() int {
	if m.Expression != nil {
		return m.Expression.ID
	}

	return 0
}

func (m *Event) StrategyID() int {
	if m.Strategy != nil {
		return m.Strategy.ID
	}

	return 0
}

func (m *Event) TemplateID() int {
	if m.Strategy != nil {
		return m.Strategy.Template.ID
	}

	return 0
}

func (m *Event) Template() *Template {
	if m.Strategy != nil {
		return m.Strategy.Template
	}

	return nil
}

func (m *Event) ActionID() int {
	if m.Expression != nil {
		return m.Expression.ActionID
	}

	return m.Strategy.Template.ActionID
}

func (m *Event) Priority() int {
	if m.Strategy != nil {
		return m.Strategy.Priority
	}
	return m.Expression.Priority
}

func (m *Event) Note() string {
	if m.Strategy != nil {
		return m.Strategy.Note
	}
	return m.Expression.Note
}

func (m *Event) Metric() string {
	if m.Strategy != nil {
		return m.Strategy.Metric
	}
	return m.Expression.Metric
}

func (m *Event) RightValue() float64 {
	if m.Strategy != nil {
		return m.Strategy.RightValue
	}
	return m.Expression.RightValue
}

func (m *Event) Operator() string {
	if m.Strategy != nil {
		return m.Strategy.Operator
	}
	return m.Expression.Operator
}

func (m *Event) Func() string {
	if m.Strategy != nil {
		return m.Strategy.Func
	}
	return m.Expression.Func
}

func (m *Event) MaxStep() int {
	if m.Strategy != nil {
		return m.Strategy.MaxStep
	}
	return m.Expression.MaxStep
}

func (m *Event) Counter() string {
	return fmt.Sprintf("%s/%s %s", m.Endpoint, m.Metric(), utils.SortedTags(m.PushedTags))
}

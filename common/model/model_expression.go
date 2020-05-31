package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

type Expression struct {
	ID         int               `json:"ID"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Note       string            `json:"note"`
	ActionID   int               `json:"actionID"`
}

func (m *Expression) String() string {
	return fmt.Sprintf(
		"<ID: %d, Metric: %s, Tags: %v, %s%s%s MaxStep: %d, P%d %s ActionId: %d>",
		m.ID,
		m.Metric,
		m.Tags,
		m.Func,
		m.Operator,
		utils.ReadableFloat(m.RightValue),
		m.MaxStep,
		m.Priority,
		m.Note,
		m.ActionID,
	)
}

type ExpressionResponse struct {
	Expressions []*Expression `json:"expressions"`
}

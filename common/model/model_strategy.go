package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

type Strategy struct {
	ID         int               `json:"ID"`
	Metric     string            `json:"metric"`
	Tags       map[string]string `json:"tags"`
	MaxStep    int               `json:"maxStep"`
	Priority   int               `json:"priority"`
	Func       string            `json:"func"`       // e.g. max(#3) all(#3)
	Operator   string            `json:"operator"`   // e.g. < !=
	RightValue float64           `json:"rightValue"` // critical value
	Note       string            `json:"note"`
	Template   *Template         `json:"template"`
}

func (m *Strategy) String() string {
	return fmt.Sprintf(
		"<ID: %d, Metric: %s, Tags: %v, %s%s%s MaxStep: %d, P%d, %s, %v>",
		m.ID,
		m.Metric,
		m.Tags,
		m.Func,
		m.Operator,
		utils.ReadableFloat(m.RightValue),
		m.MaxStep,
		m.Priority,
		m.Note,
		m.Template,
	)
}

type HostStrategy struct {
	Hostname   string     `json:"hostname"`
	Strategies []Strategy `json:"strategies"`
}

type StrategiesResponse struct {
	HostStrategies []*HostStrategy `json:"hostStrategies"`
}

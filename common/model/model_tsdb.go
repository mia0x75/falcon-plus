package model

import (
	"fmt"
	"strings"
)

type TSDBItem struct {
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func (m *TSDBItem) String() string {
	return fmt.Sprintf(
		"<Metric: %s, Tags: %v, Value: %v, TS: %d>",
		m.Metric,
		m.Tags,
		m.Value,
		m.Timestamp,
	)
}

func (m *TSDBItem) TsdbString() (s string) {
	s = fmt.Sprintf("put %s %d %.3f ", m.Metric, m.Timestamp, m.Value)

	for k, v := range m.Tags {
		key := strings.ToLower(strings.Replace(k, " ", "_", -1))
		value := strings.Replace(v, " ", "_", -1)
		s += key + "=" + value + " "
	}

	return s
}

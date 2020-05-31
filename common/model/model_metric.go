package model

import (
	"fmt"

	cu "github.com/open-falcon/falcon-plus/common/utils"
)

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (m *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint: %s, Metric: %s, Type: %s, Tags: %s, Step: %d, Time: %d, Value: %v>",
		m.Endpoint,
		m.Metric,
		m.Type,
		m.Tags,
		m.Step,
		m.Timestamp,
		m.Value,
	)
}

// Same As `MetricValue`
type JSONMetaData struct {
	Metric      string      `json:"metric"`
	Endpoint    string      `json:"endpoint"`
	Timestamp   int64       `json:"timestamp"`
	Step        int64       `json:"step"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	Tags        string      `json:"tags"`
}

func (m *JSONMetaData) String() string {
	return fmt.Sprintf(
		"<JSONMetaData Endpoint: %s, Metric: %s, Tags: %s, DsType: %s, Step: %d, Value: %v, Timestamp: %d>",
		m.Endpoint,
		m.Metric,
		m.Tags,
		m.CounterType,
		m.Step,
		m.Value,
		m.Timestamp,
	)
}

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

func (m *MetaData) String() string {
	return fmt.Sprintf(
		"<MetaData Endpoint: %s, Metric: %s, Timestamp: %d, Step: %d, Value: %f, Tags: %v>",
		m.Endpoint,
		m.Metric,
		m.Timestamp,
		m.Step,
		m.Value,
		m.Tags,
	)
}

func (m *MetaData) PK() string {
	return cu.PK(m.Endpoint, m.Metric, m.Tags)
}

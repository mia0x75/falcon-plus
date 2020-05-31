package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

type JudgeItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	JudgeType string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}

func (m *JudgeItem) String() string {
	return fmt.Sprintf("<Endpoint: %s, Metric: %s, Value: %f, Timestamp: %d, JudgeType: %s Tags: %v>",
		m.Endpoint,
		m.Metric,
		m.Value,
		m.Timestamp,
		m.JudgeType,
		m.Tags,
	)
}

func (m *JudgeItem) PrimaryKey() string {
	return utils.Md5(utils.PK(m.Endpoint, m.Metric, m.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

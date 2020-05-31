package model

import (
	"fmt"

	ttime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/common/utils"
)

type NodataItem struct {
	Counter string `json:"counter"`
	Ts      int64  `json:"ts"`
	FStatus string `json:"fstatus"`
	FTs     int64  `json:"fts"`
}

func (m *NodataItem) String() string {
	return fmt.Sprintf(
		"<NodataItem counter: %s ts: %s fecthStatus: %s fetchTs: %s>",
		m.Counter,
		ttime.FormatTs(m.Ts),
		m.FStatus,
		ttime.FormatTs(m.FTs),
	)
}

type NodataConfig struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	ObjType  string            `json:"objType"`
	Endpoint string            `json:"endpoint"`
	Metric   string            `json:"metric"`
	Tags     map[string]string `json:"tags"`
	Type     string            `json:"type"`
	Step     int64             `json:"step"`
	Mock     float64           `json:"mock"`
}

func NewNodataConfig(id int, name string, objType string, endpoint string, metric string, tags map[string]string, dstype string, step int64, mock float64) *NodataConfig {
	return &NodataConfig{id, name, objType, endpoint, metric, tags, dstype, step, mock}
}

func (m *NodataConfig) String() string {
	return fmt.Sprintf(
		"<NodataConfig id: %d, name: %s, objType: %s, endpoint: %s, metric: %s, tags: %s, type: %s, step: %d, mock: %f>",
		m.Id,
		m.Name,
		m.ObjType,
		m.Endpoint,
		m.Metric,
		utils.SortedTags(m.Tags),
		m.Type,
		m.Step,
		m.Mock,
	)
}

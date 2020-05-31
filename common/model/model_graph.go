package model

import (
	"fmt"
	"math"

	MUtils "github.com/open-falcon/falcon-plus/common/utils"
)

// DsType 即RRD中的Datasource的类型：GAUGE|COUNTER|DERIVE
type GraphItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	DsType    string            `json:"dstype"`
	Step      int               `json:"step"`
	Heartbeat int               `json:"heartbeat"`
	Min       string            `json:"min"`
	Max       string            `json:"max"`
}

func (m *GraphItem) String() string {
	return fmt.Sprintf(
		"<Endpoint: %s, Metric: %s, Tags: %v, Value: %v, TS: %d %v DsType: %s, Step: %d, Heartbeat: %d, Min: %s, Max: %s>",
		m.Endpoint,
		m.Metric,
		m.Tags,
		m.Value,
		m.Timestamp,
		MUtils.UnixTsFormat(m.Timestamp),
		m.DsType,
		m.Step,
		m.Heartbeat,
		m.Min,
		m.Max,
	)
}

func (m *GraphItem) PrimaryKey() string {
	return MUtils.PK(m.Endpoint, m.Metric, m.Tags)
}

func (t *GraphItem) Checksum() string {
	return MUtils.Checksum(t.Endpoint, t.Metric, t.Tags)
}

func (m *GraphItem) UUID() string {
	return MUtils.UUID(m.Endpoint, m.Metric, m.Tags, m.DsType, m.Step)
}

type GraphDeleteParam struct {
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Step     int    `json:"step"`
	DsType   string `json:"dstype"`
	Tags     string `json:"tags"`
}

type GraphDeleteResp struct {
}

// ConsolFun 是RRD中的概念，比如：MIN|MAX|AVERAGE
type GraphQueryParam struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
	Endpoint  string `json:"endpoint"`
	Counter   string `json:"counter"`
	Step      int    `json:"step"`
}

type GraphQueryResponse struct {
	Endpoint string     `json:"endpoint"`
	Counter  string     `json:"counter"`
	DsType   string     `json:"dstype"`
	Step     int        `json:"step"`
	Values   []*RRDData `json:"Values"` // 大写为了兼容已经再用这个api的用户
}

// 页面上已经可以看到DsType和Step了，直接带进查询条件，Graph更易处理
type GraphAccurateQueryParam struct {
	Checksum  string `json:"checksum"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
	DsType    string `json:"dsType"`
	Step      int    `json:"step"`
}

type GraphAccurateQueryResponse struct {
	Values []*RRDData `json:"values"`
}

type JSONFloat float64

func (v JSONFloat) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return []byte("null"), nil
	} else {
		return []byte(fmt.Sprintf("%f", f)), nil
	}
}

type RRDData struct {
	Timestamp int64     `json:"timestamp"`
	Value     JSONFloat `json:"value"`
}

func NewRRDData(ts int64, val float64) *RRDData {
	return &RRDData{Timestamp: ts, Value: JSONFloat(val)}
}

func (m *RRDData) String() string {
	return fmt.Sprintf(
		"<RRDData:Value: %v TS: %d %v>",
		m.Value,
		m.Timestamp,
		MUtils.UnixTsFormat(m.Timestamp),
	)
}

type GraphInfoParam struct {
	Endpoint string `json:"endpoint"`
	Counter  string `json:"counter"`
}

type GraphInfoResp struct {
	ConsolFun string `json:"consolFun"`
	Step      int    `json:"step"`
	Filename  string `json:"filename"`
}

type GraphFullyInfo struct {
	Endpoint  string `json:"endpoint"`
	Counter   string `json:"counter"`
	ConsolFun string `json:"consolFun"`
	Step      int    `json:"step"`
	Filename  string `json:"filename"`
	Addr      string `json:"addr"`
}

type GraphLastParam struct {
	Endpoint string `json:"endpoint"`
	Counter  string `json:"counter"`
}

type GraphLastResp struct {
	Endpoint string   `json:"endpoint"`
	Counter  string   `json:"counter"`
	Value    *RRDData `json:"value"`
}

package rpc

import (
	"fmt"
	"strconv"
	"time"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"github.com/open-falcon/falcon-plus/modules/transfer/proc"
	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

// Transfer TODO:
type Transfer int

// TransferResp TODO:
type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

// String TODO:
func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}

// Ping TODO:
func (s *Transfer) Ping(req cm.NullRPCRequest, resp *cm.SimpleRPCResponse) error {
	return nil
}

// Update TODO:
func (s *Transfer) Update(args []*cm.MetricValue, reply *cm.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}

// RecvMetricValues process new metric values
func RecvMetricValues(args []*cm.MetricValue, reply *cm.TransferResponse, from string) error {
	start := time.Now()
	reply.Invalid = 0

	items := []*cm.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid++
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			reply.Invalid++
			continue
		}

		if v.Type != g.COUNTER && v.Type != g.GAUGE && v.Type != g.DERIVE {
			reply.Invalid++
			continue
		}

		if v.Value == "" {
			reply.Invalid++
			continue
		}

		if v.Step <= 0 {
			reply.Invalid++
			continue
		}

		if len(v.Metric)+len(v.Tags) > 510 {
			reply.Invalid++
			continue
		}

		// Original condition: v.Timestamp <= 0 || v.Timestamp > now*2
		if abs(start.Unix()-v.Timestamp) > 5 {
			v.Timestamp = start.Unix()
		}

		fv := &cm.MetaData{
			Metric:      v.Metric,
			Endpoint:    v.Endpoint,
			Timestamp:   v.Timestamp,
			Step:        v.Step,
			CounterType: v.Type,
			Tags:        cu.DictedTagstring(v.Tags), // TODO tags键值对的个数,要做一下限制
		}

		valid := true
		var vv float64
		var err error

		switch cv := v.Value.(type) {
		case string:
			vv, err = strconv.ParseFloat(cv, 64)
			if err != nil {
				valid = false
			}
		case float64:
			vv = cv
		case int64:
			vv = float64(cv)
		default:
			valid = false
		}

		if !valid {
			reply.Invalid++
			continue
		}

		fv.Value = vv
		items = append(items, fv)
	}

	// Statistics
	cnt := int64(len(items))
	proc.RecvCnt.IncrBy(cnt)
	if from == "rpc" {
		proc.RPCRecvCnt.IncrBy(cnt)
	} else if from == "http" {
		proc.HTTPRecvCnt.IncrBy(cnt)
	}

	cfg := g.Config()

	if cfg.Graph.Enabled {
		sender.Push2GraphSendQueue(items)
	}

	if cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(items)
	}

	if cfg.TSDB.Enabled {
		sender.Push2TSDBSendQueue(items)
	}

	if cfg.Transfer.Enabled {
		sender.Push2TransferSendQueue(items)
	}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

package rpc

import (
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
	"github.com/open-falcon/falcon-plus/modules/judge/store"
)

type Judge int

func (this *Judge) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (this *Judge) Send(items []*cmodel.JudgeItem, resp *cmodel.SimpleRpcResponse) error {
	remain := g.Config().Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()
	for _, item := range items {
		exists := g.FilterMap.Exists(item.Metric)
		if !exists {
			continue
		}
		pk := item.PrimaryKey()
		store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)
	}
	return nil
}

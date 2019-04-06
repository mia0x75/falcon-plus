package sender

import (
	nlist "github.com/toolkits/container/list"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func initSendQueues() {
	cfg := g.Config()
	for node := range cfg.Judge.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	if cfg.Tsdb.Enabled {
		TsdbQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}

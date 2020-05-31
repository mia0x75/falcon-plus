package sender

import (
	nlist "github.com/toolkits/container/list"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

// initSendQueues 初始化数据缓冲
func initSendQueues() {
	cfg := g.Config()

	if cfg.Judge.Enabled {
		for node := range cfg.Judge.Cluster {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			JudgeQueues[node] = Q
		}
	}

	if cfg.Graph.Enabled {
		for node, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
				GraphQueues[node+addr] = Q
			}
		}
	}

	if cfg.TSDB.Enabled {
		TSDBQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}

	if cfg.Transfer.Enabled {
		TransferQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}

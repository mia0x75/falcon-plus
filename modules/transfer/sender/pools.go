package sender

import (
	nset "github.com/toolkits/container/set"

	cp "github.com/open-falcon/falcon-plus/common/pool"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

// initConnPools 初始化链接池
func initConnPools() {
	cfg := g.Config()

	// Judge
	if cfg.Judge.Enabled {
		judgeInstances := nset.NewStringSet()
		for _, instance := range cfg.Judge.Cluster {
			judgeInstances.Add(instance)
		}
		JudgeConnPools = cp.CreateSafeRPCConnPools(cfg.Judge.MaxConnections, cfg.Judge.MaxIdle,
			cfg.Judge.ConnectTimeout, cfg.Judge.ExecuteTimeout, judgeInstances.ToSlice())
	}

	// TSDB
	if cfg.TSDB.Enabled {
		TSDBConnPoolHelper = cp.NewTSDBConnPoolHelper(cfg.TSDB.Address, cfg.TSDB.MaxConnections, cfg.TSDB.MaxIdle, cfg.TSDB.ConnectTimeout, cfg.TSDB.ExecuteTimeout)
	}

	// Graph
	if cfg.Graph.Enabled {
		graphInstances := nset.NewSafeSet()
		for _, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				graphInstances.Add(addr)
			}
		}
		GraphConnPools = cp.CreateSafeRPCConnPools(cfg.Graph.MaxConnections, cfg.Graph.MaxIdle,
			cfg.Graph.ConnectTimeout, cfg.Graph.ExecuteTimeout, graphInstances.ToSlice())
	}

	// Transfer
	if cfg.Transfer.Enabled {
		transferInstances := nset.NewStringSet()
		for hn, instance := range cfg.Transfer.Cluster {
			TransferHostnames = append(TransferHostnames, hn)
			TransferMap[hn] = instance
			transferInstances.Add(instance)
		}
		TransferConnPools = cp.CreateSafeJSONRPCConnPools(cfg.Transfer.MaxConns, cfg.Transfer.MaxIdle,
			cfg.Transfer.ConnTimeout, cfg.Transfer.CallTimeout, transferInstances.ToSlice())
	}
}

// DestroyConnPools 销毁链接池
func DestroyConnPools() {
	cfg := g.Config()

	if cfg.Judge.Enabled {
		JudgeConnPools.Destroy()
	}

	if cfg.Graph.Enabled {
		GraphConnPools.Destroy()
	}

	if cfg.TSDB.Enabled {
		TSDBConnPoolHelper.Destroy()
	}

	if cfg.Transfer.Enabled {
		TransferConnPools.Destroy()
	}
}

package sender

import (
	nset "github.com/toolkits/container/set"

	cpools "github.com/open-falcon/falcon-plus/common/backend_pool"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func initConnPools() {
	cfg := g.Config()

	// judge
	judgeInstances := nset.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = cpools.CreateSafeRpcConnPools(cfg.Judge.MaxConnections, cfg.Judge.MaxIdle,
		cfg.Judge.ConnectTimeout, cfg.Judge.ExecuteTimeout, judgeInstances.ToSlice())

	// tsdb
	if cfg.Tsdb.Enabled {
		TsdbConnPoolHelper = cpools.NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConnections, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnectTimeout, cfg.Tsdb.ExecuteTimeout)
	}

	// graph
	graphInstances := nset.NewSafeSet()
	for _, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			graphInstances.Add(addr)
		}
	}
	GraphConnPools = cpools.CreateSafeRpcConnPools(cfg.Graph.MaxConnections, cfg.Graph.MaxIdle,
		cfg.Graph.ConnectTimeout, cfg.Graph.ExecuteTimeout, graphInstances.ToSlice())

}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	TsdbConnPoolHelper.Destroy()
}

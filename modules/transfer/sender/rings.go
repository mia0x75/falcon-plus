package sender

import (
	rings "github.com/toolkits/consistent/rings"

	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func initNodeRings() {
	cfg := g.Config()

	if cfg.Judge.Enabled {
		JudgeNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), cu.KeysOfMap(cfg.Judge.Cluster))
	}

	if cfg.Graph.Enabled {
		GraphNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), cu.KeysOfMap(cfg.Graph.Cluster))
	}
}

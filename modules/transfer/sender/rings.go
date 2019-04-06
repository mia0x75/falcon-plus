package sender

import (
	rings "github.com/toolkits/consistent/rings"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func initNodeRings() {
	cfg := g.Config()

	JudgeNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), cutils.KeysOfMap(cfg.Judge.Cluster))
	GraphNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), cutils.KeysOfMap(cfg.Graph.Cluster))
}

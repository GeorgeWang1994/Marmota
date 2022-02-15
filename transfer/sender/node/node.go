package node

import (
	"github.com/toolkits/consistent/rings"
	"marmota/pkg/utils/mmap"
	"marmota/transfer/cc"
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing *rings.ConsistentHashNodeRing
	GraphNodeRing *rings.ConsistentHashNodeRing
)

func InitNodeRings() {
	cfg := cc.Config()

	JudgeNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), mmap.KeysOfMap(cfg.Judge.Cluster))
	GraphNodeRing = rings.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), mmap.KeysOfMap(cfg.Graph.Cluster))
}

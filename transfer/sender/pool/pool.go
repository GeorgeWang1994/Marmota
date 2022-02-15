package pool

import (
	"marmota/pkg/utils/connpool"
	"marmota/transfer/cc"

	nset "github.com/toolkits/container/set"
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools    *connpool.SafeRpcConnPools
	GraphConnPools    *connpool.SafeRpcConnPools
	TransferConnPools *connpool.SafeRpcConnPools
)

// transfer的主机列表，以及主机名和地址的映射关系
// 用于随机遍历transfer
var (
	TransferMap       = make(map[string]string, 0)
	TransferHostnames = make([]string, 0)
)

// InitConnPools 初始化各个组件的连接池
func InitConnPools() {
	cfg := cc.Config()

	// judge
	judgeInstances := nset.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = connpool.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())

	// tsdb
	//if cfg.Tsdb.Enabled {
	//	TsdbConnPoolHelper = connpool.NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
	//}

	// graph
	graphInstances := nset.NewSafeSet()
	for _, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			graphInstances.Add(addr)
		}
	}
	GraphConnPools = connpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())

	// transfer
	if cfg.Transfer.Enabled {
		transferInstances := nset.NewStringSet()
		for hn, instance := range cfg.Transfer.Cluster {
			TransferHostnames = append(TransferHostnames, hn)
			TransferMap[hn] = instance
			transferInstances.Add(instance)
		}
		TransferConnPools = connpool.CreateSafeRpcConnPools(cfg.Transfer.MaxConns, cfg.Transfer.MaxIdle,
			cfg.Transfer.ConnTimeout, cfg.Transfer.CallTimeout, transferInstances.ToSlice())
	}
}

func DestroyConnPools() {
	cfg := cc.Config()

	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()

	//if cfg.Tsdb.Enabled {
	//	TsdbConnPoolHelper.Destroy()
	//}

	if cfg.Transfer.Enabled {
		TransferConnPools.Destroy()
	}
}

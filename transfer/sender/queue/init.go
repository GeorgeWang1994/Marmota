package queue

import (
	nlist "github.com/toolkits/container/list"
	"marmota/transfer/cc"
	"marmota/transfer/gg"
)

// InitSendQueues 初始化队列
func InitSendQueues() {
	cfg := cc.Config()

	for node := range cfg.Judge.Cluster {
		Q := nlist.NewSafeListLimited(gg.DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	for node, item := range cfg.Graph.ClusterList {
		for _, addr := range item.Addrs {
			Q := nlist.NewSafeListLimited(gg.DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	if cfg.Tsdb.Enabled {
		TsdbQueue = nlist.NewSafeListLimited(gg.DefaultSendQueueMaxSize)
	}

	if cfg.Transfer.Enabled {
		TransferQueue = nlist.NewSafeListLimited(gg.DefaultSendQueueMaxSize)
	}

	if cfg.Influxdb.Enabled {
		InfluxdbQueue = nlist.NewSafeListLimited(gg.DefaultSendQueueMaxSize)
	}
}

package task

import (
	"marmota/transfer/cc"
	"marmota/transfer/sender/queue"
)

// StartSendTasks TODO 添加对发送任务的控制,比如stop等
func StartSendTasks() {
	cfg := cc.Config()
	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConns
	graphConcurrent := cfg.Graph.MaxConns
	tsdbConcurrent := cfg.Tsdb.MaxConns
	transferConcurrent := cfg.Transfer.MaxConns
	influxdbConcurrent := cfg.Influxdb.MaxConns

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}

	if transferConcurrent < 1 {
		transferConcurrent = 1
	}

	if influxdbConcurrent < 1 {
		influxdbConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Judge.Cluster {
		q := queue.JudgeQueues[node]
		go forward2JudgeTask(q, node, judgeConcurrent)
	}

	//for node, nitem := range cfg.Graph.ClusterList {
	//	for _, addr := range nitem.Addrs {
	//		q := queue.GraphQueues[node+addr]
	//		go forward2GraphTask(q, node, addr, graphConcurrent)
	//	}
	//}
	//
	//if cfg.Tsdb.Enabled {
	//	go forward2TsdbTask(tsdbConcurrent)
	//}
	//
	//if cfg.Transfer.Enabled {
	//	concurrent := transferConcurrent * len(cfg.Transfer.Cluster)
	//	go forward2TransferTask(queue.TransferQueue, concurrent)
	//}
	//
	//if cfg.Influxdb.Enabled {
	//	go forward2InfluxdbTask(influxdbConcurrent)
	//}
}

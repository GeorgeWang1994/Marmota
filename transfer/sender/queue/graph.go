package queue

import (
	"fmt"
	nlist "github.com/toolkits/container/list"
	"log"
	"marmota/pkg/common/model"
	"marmota/transfer/cc"
	"marmota/transfer/gg"
	"marmota/transfer/sender/node"
	"marmota/transfer/stat"
)

var (
	GraphQueues = make(map[string]*nlist.SafeListLimited)
)

// Push2GraphSendQueue 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*model.MetaData) {
	cfg := cc.Config().Graph

	for _, item := range items {
		graphItem, err := convert2GraphItem(item)
		if err != nil {
			log.Println("E:", err)
			continue
		}
		pk := item.PK()

		//// statistics. 为了效率,放到了这里,因此只有graph是enbale时才能trace
		//stat.RecvDataTrace.Trace(pk, item)
		//stat.RecvDataFilter.Filter(pk, item.Value, item)

		node, err := node.GraphNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		cnode := cfg.ClusterList[node]
		errCnt := 0
		for _, addr := range cnode.Addrs {
			Q := GraphQueues[node+addr]
			if !Q.PushFront(graphItem) {
				errCnt += 1
			}
		}

		// statistics
		if errCnt > 0 {
			stat.SendToGraphDropCnt.Incr()
		}
	}
}

// 打到Graph的数据,要根据rrdtool的特定 来限制 step、counterType、timestamp
func convert2GraphItem(d *model.MetaData) (*model.StoreItem, error) {
	item := &model.StoreItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < gg.DefaultStep {
		item.Step = gg.DefaultStep
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == gg.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == gg.COUNTER {
		item.DsType = gg.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == gg.DERIVE {
		item.DsType = gg.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) //item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

package queue

import (
	nlist "github.com/toolkits/container/list"
	"log"
	"marmota/pkg/common/model"
	"marmota/transfer/gg"
	"marmota/transfer/sender/node"
	"marmota/transfer/stat"
)

var (
	JudgeQueues = make(map[string]*nlist.SafeListLimited)
)

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}

// Push2JudgeSendQueue 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*model.MetaData) {
	for _, item := range items {
		pk := item.PK()
		n, err := node.JudgeNodeRing.GetNode(pk)
		if err != nil {
			log.Println("E:", err)
			continue
		}

		// align ts ???
		step := int(item.Step)
		if step < gg.DefaultMinStep {
			step = gg.DefaultMinStep
		}
		ts := alignTs(item.Timestamp, int64(step))

		judgeItem := &model.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := JudgeQueues[n]
		isSuccess := Q.PushFront(judgeItem)

		// statistics
		if !isSuccess {
			stat.SendToJudgeDropCnt.Incr()
		}
	}
}

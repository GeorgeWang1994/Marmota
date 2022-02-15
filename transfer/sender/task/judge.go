package task

import (
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
	"log"
	"marmota/pkg/common/model"
	"marmota/transfer/cc"
	"marmota/transfer/sender/pool"
	"marmota/transfer/stat"
	"time"
)

const DefaultSendTaskSleepInterval = time.Millisecond * 50 // 默认睡眠间隔为50ms

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *list.SafeListLimited, node string, concurrent int) {
	batch := cc.Config().Judge.Batch // 一次发送,最多batch条数据
	addr := cc.Config().Judge.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems := make([]*model.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*model.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, judgeItems []*model.JudgeItem, count int) {
			defer sema.Release()

			resp := &model.RpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = pool.JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Printf("send judge %s:%s fail: %v", node, addr, err)
				stat.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				stat.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

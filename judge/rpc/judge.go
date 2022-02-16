package rpc

import (
	"marmota/judge/cc"
	"marmota/judge/gg"
	"marmota/judge/store"
	"marmota/pkg/common/model"
	"time"
)

type Judge int

func (j *Judge) Ping(req model.NullRpcRequest, resp *model.RpcResponse) error {
	return nil
}

// Send 发送数据给存储组件
func (j *Judge) Send(items []*model.JudgeItem, resp *model.RpcResponse) error {
	remain := cc.Config().Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()
	for _, item := range items {
		exists := gg.FilterMap.Exists(item.Metric)
		if !exists {
			continue
		}
		pk := item.PK()
		//加入到链表中
		store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)
	}
	return nil
}

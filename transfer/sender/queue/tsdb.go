package queue

import (
	"marmota/pkg/common/model"
	"marmota/transfer/stat"

	nlist "github.com/toolkits/container/list"
)

var (
	TsdbQueue *nlist.SafeListLimited
)

// Push2TsdbSendQueue 将原始数据入到tsdb发送缓存队列
func Push2TsdbSendQueue(items []*model.MetaData) {
	for _, item := range items {
		tsdbItem := convert2TSDBItem(item)
		isSuccess := TsdbQueue.PushFront(tsdbItem)

		if !isSuccess {
			stat.SendToTsdbDropCnt.Incr()
		}
	}
}

// 转化为tsdb格式
func convert2TSDBItem(d *model.MetaData) *model.TSDBItem {
	t := model.TSDBItem{Tags: make(map[string]string)}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}

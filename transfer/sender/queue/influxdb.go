package queue

import (
	nlist "github.com/toolkits/container/list"
	"marmota/pkg/common/model"
	"marmota/transfer/stat"
)

var (
	InfluxdbQueue *nlist.SafeListLimited
)

// Push2InfluxdbSendQueue 将原始数据插入到influxdb缓存队列
func Push2InfluxdbSendQueue(items []*model.MetaData) {
	for _, item := range items {
		influxdbItem := convert2InfluxdbItem(item)
		isSuccess := InfluxdbQueue.PushFront(influxdbItem)

		if !isSuccess {
			stat.SendToInfluxdbDropCnt.Incr()
		}
	}
}

func convert2InfluxdbItem(d *model.MetaData) *model.InfluxdbItem {
	t := model.InfluxdbItem{Tags: make(map[string]string), Fileds: make(map[string]interface{})}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Measurement = d.Metric
	t.Fileds["value"] = d.Value
	t.Timestamp = d.Timestamp

	return &t
}

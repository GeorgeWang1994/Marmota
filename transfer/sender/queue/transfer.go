package queue

import (
	nlist "github.com/toolkits/container/list"
	"marmota/pkg/common/model"
	"marmota/transfer/stat"
)

var (
	TransferQueue *nlist.SafeListLimited
)

func Push2TransferSendQueue(items []*model.MetaData) {
	for _, item := range items {
		isSuccess := TransferQueue.PushFront(item)

		if !isSuccess {
			stat.SendToTransferDropCnt.Incr()
		}
	}
}

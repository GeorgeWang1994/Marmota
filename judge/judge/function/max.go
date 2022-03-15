package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type MaxFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f MaxFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	// 找到最大值
	max := vs[0].Value
	for i := 1; i < f.Limit; i++ {
		if max < vs[i].Value {
			max = vs[i].Value
		}
	}

	// 拿过去n个中的最大值和策略中的数据进行比较
	leftValue = max
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

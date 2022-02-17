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

	max := vs[0].Value
	for i := 1; i < f.Limit; i++ {
		if max < vs[i].Value {
			max = vs[i].Value
		}
	}

	leftValue = max
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type SumFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f SumFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < f.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

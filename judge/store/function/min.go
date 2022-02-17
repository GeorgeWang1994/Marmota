package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type MinFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f MinFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	min := vs[0].Value
	for i := 1; i < f.Limit; i++ {
		if min > vs[i].Value {
			min = vs[i].Value
		}
	}

	leftValue = min
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

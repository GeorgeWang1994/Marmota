package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type AllFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f AllFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	for i := 0; i < f.Limit; i++ {
		isTriggered = checkIsTriggered(vs[i].Value, f.Operator, f.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Value
	return
}

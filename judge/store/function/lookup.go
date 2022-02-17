package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type LookupFunction struct {
	Function
	Num        int
	Limit      int
	Operator   string
	RightValue float64
}

func (f LookupFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	leftValue = vs[0].Value

	for n, i := 0, 0; i < f.Limit; i++ {
		if checkIsTriggered(vs[i].Value, f.Operator, f.RightValue) {
			n++
			if n == f.Num {
				isTriggered = true
				return
			}
		}
	}

	return
}

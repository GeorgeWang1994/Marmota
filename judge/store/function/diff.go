package function

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

type DiffFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

// Compute 只要有一个点的diff触发阈值，就报警
func (f DiffFunction) Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	// 此处this.Limit要+1，因为通常说diff(#3)，是当前点与历史的3个点相比较
	// 然而最新点已经在linkedlist的第一个位置，所以……
	vs, isEnough = L.HistoryData(f.Limit + 1)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	first := vs[0].Value

	isTriggered = false
	for i := 1; i < f.Limit+1; i++ {
		// diff是当前值减去历史值
		leftValue = first - vs[i].Value
		isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
		if isTriggered {
			break
		}
	}

	return
}

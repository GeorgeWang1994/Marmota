package judge

import (
	"marmota/judge/store"
	"marmota/pkg/common/model"
)

func Judge(L *store.SafeLinkedList, firstItem *model.JudgeItem, now int64) {
	CheckStrategy(L, firstItem, now)
	//CheckExpression(L, firstItem, now)
}

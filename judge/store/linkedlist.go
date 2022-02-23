package store

import (
	"container/list"
	"marmota/pkg/common/model"
	"sync"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (s *SafeLinkedList) ToSlice() []*model.JudgeItem {
	s.RLock()
	defer s.RUnlock()
	sz := s.L.Len()
	if sz == 0 {
		return []*model.JudgeItem{}
	}

	ret := make([]*model.JudgeItem, 0, sz)
	for e := s.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*model.JudgeItem))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (s *SafeLinkedList) HistoryData(limit int) ([]*model.HistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*model.HistoryData{}, false
	}

	size := s.Len()
	if size == 0 {
		return []*model.HistoryData{}, false
	}

	firstElement := s.Front()
	firstItem := firstElement.Value.(*model.JudgeItem)

	var vs []*model.HistoryData
	isEnough := true

	judgeType := firstItem.JudgeType[0]
	if judgeType == 'G' || judgeType == 'g' {
		if size < limit {
			// 有多少获取多少
			limit = size
			isEnough = false
		}
		vs = make([]*model.HistoryData, limit)
		vs[0] = &model.HistoryData{Timestamp: firstItem.Timestamp, Value: firstItem.Value}
		i := 1
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			vs[i] = &model.HistoryData{
				Timestamp: nextElement.Value.(*model.JudgeItem).Timestamp,
				Value:     nextElement.Value.(*model.JudgeItem).Value,
			}
			i++
			currentElement = nextElement
		}
	} else {
		if size < limit+1 {
			isEnough = false
			limit = size - 1
		}

		vs = make([]*model.HistoryData, limit)

		i := 0
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			diffVal := currentElement.Value.(*model.JudgeItem).Value - nextElement.Value.(*model.JudgeItem).Value
			diffTs := currentElement.Value.(*model.JudgeItem).Timestamp - nextElement.Value.(*model.JudgeItem).Timestamp
			vs[i] = &model.HistoryData{
				Timestamp: currentElement.Value.(*model.JudgeItem).Timestamp,
				Value:     diffVal / float64(diffTs),
			}
			i++
			currentElement = nextElement
		}
	}

	return vs, isEnough
}

func (s *SafeLinkedList) PushFront(v interface{}) *list.Element {
	s.Lock()
	defer s.Unlock()
	return s.L.PushFront(v)
}

// @return needJudge 如果是false不需要做judge，因为新上来的数据不合法
func (s *SafeLinkedList) PushFrontAndMaintain(v *model.JudgeItem, maxCount int) bool {
	s.Lock()
	defer s.Unlock()

	sz := s.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.Timestamp <= s.L.Front().Value.(*model.JudgeItem).Timestamp || v.Timestamp <= 0 {
			return false
		}
	}

	s.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		s.L.Remove(s.L.Back())
	}

	return true
}

func (s *SafeLinkedList) Front() *list.Element {
	s.RLock()
	defer s.RUnlock()
	return s.L.Front()
}

func (s *SafeLinkedList) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.L.Len()
}

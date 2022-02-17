package store

import (
	"container/list"
	"marmota/judge/store/judge"
	"marmota/pkg/common/model"
	"sync"
)

type JudgeItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList
}

func NewJudgeItemMap() *JudgeItemMap {
	return &JudgeItemMap{M: make(map[string]*SafeLinkedList)}
}

func (j *JudgeItemMap) Get(key string) (*SafeLinkedList, bool) {
	j.RLock()
	defer j.RUnlock()
	val, ok := j.M[key]
	return val, ok
}

func (j *JudgeItemMap) Set(key string, val *SafeLinkedList) {
	j.Lock()
	defer j.Unlock()
	j.M[key] = val
}

func (j *JudgeItemMap) Len() int {
	j.RLock()
	defer j.RUnlock()
	return len(j.M)
}

func (j *JudgeItemMap) Delete(key string) {
	j.Lock()
	defer j.Unlock()
	delete(j.M, key)
}

func (j *JudgeItemMap) BatchDelete(keys []string) {
	count := len(keys)
	if count == 0 {
		return
	}

	j.Lock()
	defer j.Unlock()
	for i := 0; i < count; i++ {
		delete(j.M, keys[i])
	}
}

// CleanStale 找出所有小于before的JudgeItem，左后清理掉
func (j *JudgeItemMap) CleanStale(before int64) {
	var keys []string

	j.RLock()
	for key, L := range j.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*model.JudgeItem).Timestamp < before {
			keys = append(keys, key)
		}
	}
	j.RUnlock()

	j.BatchDelete(keys)
}

func (j *JudgeItemMap) PushFrontAndMaintain(key string, val *model.JudgeItem, maxCount int, now int64) {
	if linkedList, exists := j.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			judge.Judge(linkedList, val, now)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
		j.Set(key, safeList)
		judge.Judge(safeList, val, now)
	}
}

// HistoryBigMap 这是个线程不安全的大Map，需要提前初始化好
var HistoryBigMap = make(map[string]*JudgeItemMap)

func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewJudgeItemMap()
		}
	}
}

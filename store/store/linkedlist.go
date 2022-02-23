package store

import (
	"container/list"
	"marmota/pkg/common/model"
	"sync"
)

type SafeLinkedList struct {
	sync.RWMutex
	Flag uint32
	L    *list.List
}

// 新创建SafeLinkedList容器
func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

func (s *SafeLinkedList) PushFront(v interface{}) *list.Element {
	s.Lock()
	defer s.Unlock()
	return s.L.PushFront(v)
}

func (s *SafeLinkedList) Front() *list.Element {
	s.RLock()
	defer s.RUnlock()
	return s.L.Front()
}

func (s *SafeLinkedList) PopBack() *list.Element {
	s.Lock()
	defer s.Unlock()

	back := s.L.Back()
	if back != nil {
		s.L.Remove(back)
	}

	return back
}

func (s *SafeLinkedList) Back() *list.Element {
	s.Lock()
	defer s.Unlock()

	return s.L.Back()
}

func (s *SafeLinkedList) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.L.Len()
}

// remain参数表示要给linkedlist中留几个元素
// 在cron中刷磁盘的时候要留一个，用于创建数据库索引
// 在程序退出的时候要一个不留的全部刷到磁盘
func (s *SafeLinkedList) PopAll() []*model.StoreItem {
	s.Lock()
	defer s.Unlock()

	size := s.L.Len()
	if size <= 0 {
		return []*model.StoreItem{}
	}

	ret := make([]*model.StoreItem, 0, size)

	for i := 0; i < size; i++ {
		item := s.L.Back()
		ret = append(ret, item.Value.(*model.StoreItem))
		s.L.Remove(item)
	}

	return ret
}

//restore PushAll
func (s *SafeLinkedList) PushAll(items []*model.StoreItem) {
	s.Lock()
	defer s.Unlock()

	size := len(items)
	if size > 0 {
		for i := size - 1; i >= 0; i-- {
			s.L.PushBack(items[i])
		}
	}
}

//return为倒叙的?
func (s *SafeLinkedList) FetchAll() ([]*model.StoreItem, uint32) {
	s.Lock()
	defer s.Unlock()
	count := s.L.Len()
	ret := make([]*model.StoreItem, 0, count)

	p := s.L.Back()
	for p != nil {
		ret = append(ret, p.Value.(*model.StoreItem))
		p = p.Prev()
	}

	return ret, s.Flag
}

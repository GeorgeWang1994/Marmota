package gg

import (
	"marmota/pkg/common/model"
	"sync"
)

type SafeEventMap struct {
	sync.RWMutex
	M map[string]*model.Event
}

var (
	LastEvents = &SafeEventMap{M: make(map[string]*model.Event)} // 记录上次告警的事件
)

func (s *SafeEventMap) Get(key string) (*model.Event, bool) {
	s.RLock()
	defer s.RUnlock()
	event, exists := s.M[key]
	return event, exists
}

func (s *SafeEventMap) Set(key string, event *model.Event) {
	s.Lock()
	defer s.Unlock()
	s.M[key] = event
}

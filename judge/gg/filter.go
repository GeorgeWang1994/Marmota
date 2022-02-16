package gg

import "sync"

type SafeFilterMap struct {
	sync.RWMutex
	M map[string]string
}

var (
	FilterMap = &SafeFilterMap{M: make(map[string]string)}
)

func (s *SafeFilterMap) ReInit(m map[string]string) {
	s.Lock()
	defer s.Unlock()
	s.M = m
}

func (s *SafeFilterMap) Exists(key string) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.M[key]; ok {
		return true
	}
	return false
}

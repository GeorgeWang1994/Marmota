package gg

import (
	"marmota/pkg/common/model"
	"sync"
)

type SafeStrategyMap struct {
	sync.RWMutex
	// endpoint:metric => [strategy1, strategy2 ...]
	M map[string][]model.Strategy
}

var (
	StrategyMap = &SafeStrategyMap{M: make(map[string][]model.Strategy)} // 缓存终端和策略的映射关系
)

func (s *SafeStrategyMap) ReInit(m map[string][]model.Strategy) {
	s.Lock()
	defer s.Unlock()
	s.M = m
}

func (s *SafeStrategyMap) Get() map[string][]model.Strategy {
	s.RLock()
	defer s.RUnlock()
	return s.M
}

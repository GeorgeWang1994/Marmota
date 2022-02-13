package cache

import (
	"marmota/pivas/db"
	"marmota/pkg/common/model"
	"sync"
)

type SafeExpressionCache struct {
	sync.RWMutex
	L []*model.Expression
}

var ExpressionCache = &SafeExpressionCache{}

func (s *SafeExpressionCache) Get() []*model.Expression {
	s.RLock()
	defer s.RUnlock()
	return s.L
}

func (s *SafeExpressionCache) Init() {
	es, err := db.QueryExpressions()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.L = es
}


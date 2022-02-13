package cache

import (
	"marmota/pivas/db"
	"marmota/pkg/common/model"
	"sync"
)

// SafeGroupTemplates 一个HostGroup对应多个Template
type SafeGroupTemplates struct {
	sync.RWMutex
	M map[int][]int
}

var GroupTemplates = &SafeGroupTemplates{M: make(map[int][]int)}

func (s *SafeGroupTemplates) GetTemplateIds(gid int) ([]int, bool) {
	s.RLock()
	defer s.RUnlock()
	templateIds, exists := s.M[gid]
	return templateIds, exists
}

func (s *SafeGroupTemplates) Init() {
	m, err := db.QueryGroupTemplates()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = m
}

type SafeTemplateCache struct {
	sync.RWMutex
	M map[int]*model.Template
}

var TemplateCache = &SafeTemplateCache{M: make(map[int]*model.Template)}

func (s *SafeTemplateCache) GetMap() map[int]*model.Template {
	s.RLock()
	defer s.RUnlock()
	return s.M
}

func (s *SafeTemplateCache) Init() {
	ts, err := db.QueryTemplates()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = ts
}

type SafeHostTemplateIds struct {
	sync.RWMutex
	M map[int][]int
}

var HostTemplateIds = &SafeHostTemplateIds{M: make(map[int][]int)}

func (s *SafeHostTemplateIds) GetMap() map[int][]int {
	s.RLock()
	defer s.RUnlock()
	return s.M
}

func (s *SafeHostTemplateIds) Init() {
	m, err := db.QueryHostTemplateIds()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = m
}

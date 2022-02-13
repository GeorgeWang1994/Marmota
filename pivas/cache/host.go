package cache

import (
	"marmota/pivas/db"
	"marmota/pkg/common/model"
	"sync"
)

// SafeHostMap 每次心跳的时候agent把hostname汇报上来，经常要知道这个机器的hostid，把此信息缓存
// key: hostname value: hostid
type SafeHostMap struct {
	sync.RWMutex
	M map[string]int
}

var HostMap = &SafeHostMap{M: make(map[string]int)}

func (s *SafeHostMap) GetID(hostname string) (int, bool) {
	s.RLock()
	defer s.RUnlock()
	id, exists := s.M[hostname]
	return id, exists
}

func (s *SafeHostMap) Init() {
	m, err := db.QueryHosts()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = m
}

type SafeMonitoredHosts struct {
	sync.RWMutex
	M map[int]*model.Host
}

var MonitoredHosts = &SafeMonitoredHosts{M: make(map[int]*model.Host)}

func (s *SafeMonitoredHosts) Get() map[int]*model.Host {
	s.RLock()
	defer s.RUnlock()
	return s.M
}

func (s *SafeMonitoredHosts) Init() {
	m, err := db.QueryMonitoredHosts()
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = m
}

package cache

import (
	"marmota/pivas/db"
	"marmota/pkg/common/model"
	"sync"
	"time"
)

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

type SafeAgents struct {
	sync.RWMutex
	M map[string]*model.AgentUpdateInfo
}

var Agents = NewSafeAgents()

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*model.AgentUpdateInfo)}
}

func (s *SafeAgents) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := s.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		db.UpdateAgent(val)
	}

	// 更新hbs 时间
	s.Lock()
	s.M[req.Hostname] = val
	s.Unlock()
}

func (s *SafeAgents) Get(hostname string) (*model.AgentUpdateInfo, bool) {
	s.RLock()
	defer s.RUnlock()
	val, exists := s.M[hostname]
	return val, exists
}

func (s *SafeAgents) Delete(hostname string) {
	s.Lock()
	defer s.Unlock()
	delete(s.M, hostname)
}

func (s *SafeAgents) Keys() []string {
	s.RLock()
	defer s.RUnlock()
	count := len(s.M)
	keys := make([]string, count)
	i := 0
	for hostname := range s.M {
		keys[i] = hostname
		i++
	}
	return keys
}

func DeleteStaleAgents() {
	duration := time.Hour * time.Duration(24)
	for {
		time.Sleep(duration)
		deleteStaleAgents()
	}
}

func deleteStaleAgents() {
	// 一天都没有心跳的Agent，从内存中干掉
	before := time.Now().Unix() - 3600*24
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			Agents.Delete(curr.ReportRequest.Hostname)
		}
	}
}


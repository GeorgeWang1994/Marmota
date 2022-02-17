package cache

import (
	"fmt"
	"github.com/toolkits/container/set"
	"log"
	"marmota/pivas/db"
	"marmota/pkg/common/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SafeStrategies struct {
	sync.RWMutex
	M map[int]*model.Strategy
}

var Strategies = &SafeStrategies{M: make(map[int]*model.Strategy)}

func (s *SafeStrategies) GetMap() map[int]*model.Strategy {
	s.RLock()
	defer s.RUnlock()
	return s.M
}

func (s *SafeStrategies) Init(tpls map[int]*model.Template) {
	m, err := db.QueryStrategies(tpls)
	if err != nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.M = m
}

// QueryStrategies 获取所有的Strategy列表
func QueryStrategies(tpls map[int]*model.Template) (map[int]*model.Strategy, error) {
	ret := make(map[int]*model.Strategy)

	if tpls == nil || len(tpls) == 0 {
		return ret, fmt.Errorf("illegal argument")
	}

	now := time.Now().Format("15:04")
	sql := fmt.Sprintf(
		"select %s from strategy as s where (s.run_begin='' and s.run_end='') "+
			"or (s.run_begin <= '%s' and s.run_end >= '%s')"+
			"or (s.run_begin > s.run_end and !(s.run_begin > '%s' and s.run_end < '%s'))",
		"s.id, s.metric, s.tags, s.function, s.op, s.right_value, s.max_step, s.priority, s.note, s.tpl_id",
		now,
		now,
		now,
		now,
	)

	rows, err := db.DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		s := model.Strategy{}
		var tags string
		var tid int
		err = rows.Scan(&s.Id, &s.Metric, &tags, &s.Func, &s.Operator, &s.RightValue, &s.MaxStep, &s.Priority, &s.Note, &tid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		tt := make(map[string]string)

		if tags != "" {
			arr := strings.Split(tags, ",")
			for _, tag := range arr {
				kv := strings.SplitN(tag, "=", 2)
				if len(kv) != 2 {
					continue
				}
				tt[kv[0]] = kv[1]
			}
		}

		s.Tags = tt
		s.Tpl = tpls[tid]
		if s.Tpl == nil {
			log.Printf("WARN: tpl is nil. strategy id=%d, tpl id=%d", s.Id, tid)
			// 如果Strategy没有对应的Tpl，那就没有action，就没法报警，无需往后传递了
			continue
		}

		ret[s.Id] = &s
	}

	return ret, nil
}

func GetBuiltinMetrics(hostname string) ([]*model.AgentMetric, error) {
	ret := []*model.AgentMetric{}
	hid, exists := HostMap.GetID(hostname)
	if !exists {
		return ret, nil
	}

	gids, exists := HostGroupsMap.GetGroupIds(hid)
	if !exists {
		return ret, nil
	}

	// 根据gids，获取绑定的所有tids
	tidSet := set.NewIntSet()
	for _, gid := range gids {
		tids, exists := GroupTemplates.GetTemplateIds(gid)
		if !exists {
			continue
		}

		for _, tid := range tids {
			tidSet.Add(tid)
		}
	}

	tidSlice := tidSet.ToSlice()
	if len(tidSlice) == 0 {
		return ret, nil
	}

	// 继续寻找这些tid的ParentId
	allTpls := TemplateCache.GetMap()
	for _, tid := range tidSlice {
		pids := ParentStrategyIds(allTpls, tid)
		for _, pid := range pids {
			tidSet.Add(pid)
		}
	}

	// 终于得到了最终的tid列表
	tidSlice = tidSet.ToSlice()

	// 把tid列表用逗号拼接在一起
	count := len(tidSlice)
	tidStrArr := make([]string, count)
	for i := 0; i < count; i++ {
		tidStrArr[i] = strconv.Itoa(tidSlice[i])
	}

	return db.QueryAgentMetrics(strings.Join(tidStrArr, ","))
}

func ParentStrategyIds(allTpls map[int]*model.Template, tid int) (ret []int) {
	depth := 0
	for {
		if tid <= 0 {
			break
		}

		if t, exists := allTpls[tid]; exists {
			ret = append(ret, tid)
			tid = t.ParentId
		} else {
			break
		}

		depth++
		if depth == 10 {
			log.Println("[ERROR] template inherit cycle. id:", tid)
			return []int{}
		}
	}

	sz := len(ret)
	if sz <= 1 {
		return
	}

	desc := make([]int, sz)
	for i, item := range ret {
		j := sz - i - 1
		desc[j] = item
	}

	return desc
}


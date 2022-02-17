package strategy

import (
	"encoding/json"
	"fmt"
	"log"
	"marmota/judge/cc"
	"marmota/judge/gg"
	"marmota/pkg/common/model"
)

// 定时从hbs获取所有的策略，并且更新到缓存
func syncStrategies() {
	var strategiesResponse model.StrategiesResponse
	err := gg.HBSRpcClient().Call("Hbs.GetStrategies", model.NullRpcRequest{}, &strategiesResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetStrategies:", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

// 重建本地缓存StrategyMap
func rebuildStrategyMap(strategiesResponse *model.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]model.Strategy)
	for _, hs := range strategiesResponse.HostStrategies {
		hostname := hs.Hostname
		if cc.Config().Debug && hostname == cc.Config().DebugHost {
			log.Println(hostname, "strategies:")
			bs, _ := json.Marshal(hs.Strategies)
			fmt.Println(string(bs))
		}
		for _, strategy := range hs.Strategies {
			key := fmt.Sprintf("%s/%s", hostname, strategy.Metric)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], strategy)
			} else {
				m[key] = []model.Strategy{strategy}
			}
		}
	}

	gg.StrategyMap.ReInit(m)
}

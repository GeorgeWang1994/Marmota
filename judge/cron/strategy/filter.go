package strategy

import "marmota/judge/gg"

func syncFilter() {
	m := make(map[string]string)

	//M map[string][]model.Strategy
	strategyMap := gg.StrategyMap.Get()
	for _, strategies := range strategyMap {
		for _, strategy := range strategies {
			m[strategy.Metric] = strategy.Metric
		}
	}

	gg.FilterMap.ReInit(m)
}

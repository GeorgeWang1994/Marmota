package strategy

import (
	"marmota/judge/cc"
	"time"
)

// SyncStrategies 更新策略信息
func SyncStrategies() {
	duration := time.Duration(cc.Config().Hbs.Interval) * time.Second
	for {
		syncStrategies()
		//syncExpression()
		syncFilter()
		time.Sleep(duration)
	}
}

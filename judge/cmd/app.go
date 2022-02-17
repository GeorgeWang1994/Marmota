package main

import (
	"flag"
	"marmota/judge/cc"
	"marmota/judge/cron"
	"marmota/judge/cron/strategy"
	"marmota/judge/gg"
	"marmota/judge/rpc"
	"marmota/judge/store"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")

	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	gg.InitRedisConnPool()
	store.InitHistoryBigMap()
	go strategy.SyncStrategies()
	go cron.CleanStale()

	go rpc.Start()

	return nil
}

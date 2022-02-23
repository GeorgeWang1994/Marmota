package main

import (
	"flag"
	"marmota/alarm/cc"
	"marmota/alarm/cron/event/combine"
	"marmota/alarm/cron/event/consume"
	"marmota/alarm/gg"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")
	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	gg.InitRedisConnPool()

	go consume.ReadHighEvent()
	go consume.ReadLowEvent()
	go consume.ConsumeIM()
	go combine.CombineIM()
	// todo: clear expired event
	return nil
}

package main

import (
	"flag"
	"marmota/alarm/cc"
	"marmota/alarm/gg"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")
	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	gg.InitRedisConnPool()
	return nil
}

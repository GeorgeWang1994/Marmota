package main

import (
	"flag"
	"marmota/judge/cc"
	"marmota/judge/gg"
	"marmota/judge/rpc"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")

	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	gg.InitRedisConnPool()
	go rpc.Start()

	return nil
}

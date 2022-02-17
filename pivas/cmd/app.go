package main

import (
	"flag"
	"marmota/pivas/cc"
	"marmota/pivas/rpc"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")
	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	go rpc.Start()
	return nil
}

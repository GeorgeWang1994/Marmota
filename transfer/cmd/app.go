package main

import (
	"flag"
	"marmota/transfer/cc"
	"marmota/transfer/rpc"
	"marmota/transfer/sender"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")
	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	rpc.Start()
	sender.Start()

	return nil
}

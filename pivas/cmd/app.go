package main

import (
	"marmota/pivas/cc"
	"marmota/pivas/rpc"
)

func initApp() error {
	err := cc.ParseConfig("")
	if err != nil {
		return err
	}

	go rpc.Start()
	return nil
}

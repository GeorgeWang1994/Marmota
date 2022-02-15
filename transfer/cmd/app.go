package cmd

import (
	"marmota/transfer/cc"
	"marmota/transfer/rpc"
	"marmota/transfer/sender"
)

func initApp() error {
	err := cc.ParseConfig("")
	if err != nil {
		return err
	}

	rpc.Start()
	sender.Start()

	return nil
}

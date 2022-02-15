package main

import (
	"marmota/pivas/cc"
)

func initApp() error {
	err := cc.ParseConfig("")
	if err != nil {
		return err
	}

	return nil
}

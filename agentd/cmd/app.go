package main

import (
	"marmota/agentd/cc"
	"marmota/agentd/cron"
)

func initApp() error {
	err := cc.ParseConfig("")
	if err != nil {
		return err
	}

	cron.ReportStatus()
	return nil
}

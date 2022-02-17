package main

import (
	"flag"
	"marmota/agentd/cc"
	"marmota/agentd/cron"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")
	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	cron.ReportStatus()
	cron.SyncBuiltinMetrics()
	cron.SyncMinePlugins()
	return nil
}

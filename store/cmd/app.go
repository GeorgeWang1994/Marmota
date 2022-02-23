package main

import (
	"flag"
	"marmota/store/cc"
	"marmota/store/db"
	"marmota/store/index"
	"marmota/store/rpc"
	"marmota/store/rrdtool"
)

func initApp() error {
	cfg := flag.String("c", "cfg.json", "configuration file")

	err := cc.ParseConfig(*cfg)
	if err != nil {
		return err
	}

	// 初始化数据库
	db.InitDB()
	// rrdtool init
	rrdtool.InitChannel()
	// rrdtool before api for disable loopback connection
	rrdtool.Start()
	// start api
	go rpc.Start()
	// start indexing
	// index
	index.Start()

	return nil
}

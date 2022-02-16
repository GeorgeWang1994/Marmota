package gg

import (
	"marmota/judge/cc"
	"marmota/judge/rpc"
	"time"
)

var (
	hbsRpcClient *rpc.ConnRPCClient
)

func HBSRpcClient() *rpc.ConnRPCClient {
	if hbsRpcClient != nil {
		return hbsRpcClient
	}
	hbsRpcClient = rpc.NewRpcClient(
		cc.Config().Hbs.Servers,
		time.Duration(cc.Config().Hbs.Timeout)*time.Second,
	)
	return hbsRpcClient
}

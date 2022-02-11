package gg

import "marmota/agentd/rpc"

var (
	HbsRpcClient *rpc.Client
)

func RpcClient() *rpc.Client {
	if HbsRpcClient != nil {
		return HbsRpcClient
	}
	HbsRpcClient = rpc.NewRpcClient()
	return HbsRpcClient
}

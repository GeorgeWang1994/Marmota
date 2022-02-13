package gg

import (
	"marmota/agentd/cc"
	"marmota/agentd/rpc"
	"sync"
	"time"
)

var (
	hbsRpcClient        *rpc.ConnRPCClient
	TransferClientsLock *sync.RWMutex                 = new(sync.RWMutex)
	TransferClients     map[string]*rpc.ConnRPCClient = map[string]*rpc.ConnRPCClient{}
)

func HBSRpcClient() *rpc.ConnRPCClient {
	if hbsRpcClient != nil {
		return hbsRpcClient
	}
	hbsRpcClient = rpc.NewRpcClient(
		cc.Config().Heartbeat.Addr,
		time.Duration(cc.Config().Heartbeat.Timeout)*time.Second,
	)
	return hbsRpcClient
}

func InitTransferClient(addr string) *rpc.ConnRPCClient {
	var c *rpc.ConnRPCClient = rpc.NewRpcClient(
		addr, time.Duration(cc.Config().Transfer.Timeout)*time.Millisecond,
	)
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = c

	return c
}

func TransferClient(addr string) *rpc.ConnRPCClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()

	if c, ok := TransferClients[addr]; ok {
		return c
	}
	return nil
}

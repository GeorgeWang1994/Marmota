package rpc

import (
	"errors"
	"log"
	"math"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/hashicorp/net-rpc-msgpackrpc"
)

type ConnRPCClient struct {
	sync.Mutex
	RpcAddr   string
	rpcClient *rpc.Client
	Timeout   time.Duration
}

func (c *ConnRPCClient) serverConn() error {
	if c.rpcClient != nil {
		return nil
	}

	var retry = 1

	for {
		if c.rpcClient != nil {
			return nil
		}

		conn, err := net.DialTimeout("tcp", c.RpcAddr, c.Timeout)
		if err != nil {
			log.Printf("dial %s fail: %v", c.RpcAddr, err)
			if retry > 3 {
				return err
			}
			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			continue
		}
		c.rpcClient = msgpackrpc.NewClient(conn)
		return err
	}
}

func (c *ConnRPCClient) close() {
	if c.rpcClient != nil {
		_ = c.rpcClient.Close()
		c.rpcClient = nil
	}
}

func (c *ConnRPCClient) Call(method string, args interface{}, reply interface{}) error {

	c.Lock()
	defer c.Unlock()

	err := c.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10) * time.Second
	done := make(chan error, 1)

	go func() {
		err := c.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", c.rpcClient, c.RpcAddr)
		c.close()
		return errors.New(c.RpcAddr + " rpc call timeout")
	case err := <-done:
		if err != nil {
			c.close()
			return err
		}
	}

	return nil
}

func NewRpcClient(addr string, timeout time.Duration) *ConnRPCClient {
	return &ConnRPCClient{
		RpcAddr: addr,
		Timeout: timeout,
	}
}

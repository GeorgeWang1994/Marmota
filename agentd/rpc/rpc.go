package rpc

import (
	"errors"
	"log"
	"marmota/agentd/cc"
	"math"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/hashicorp/net-rpc-msgpackrpc"
)

type Client struct {
	sync.Mutex
	rpcAddr   string
	rpcClient *rpc.Client
	Timeout   time.Duration
}

func (c *Client) serverConn() error {
	if c.rpcClient != nil {
		return nil
	}

	var retry = 1

	for {
		if c.rpcClient != nil {
			return nil
		}

		conn, err := net.DialTimeout("tcp", c.rpcAddr, c.Timeout)
		if err != nil {
			log.Printf("dial %s fail: %v", c.rpcAddr, err)
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

func (c *Client) close() {
	if c.rpcClient != nil {
		_ = c.rpcClient.Close()
		c.rpcClient = nil
	}
}

func (c *Client) Call(method string, args interface{}, reply interface{}) error {

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
		log.Printf("[WARN] rpc call timeout %v => %v", c.rpcClient, c.rpcAddr)
		c.close()
		return errors.New(c.rpcAddr + " rpc call timeout")
	case err := <-done:
		if err != nil {
			c.close()
			return err
		}
	}

	return nil
}

func NewRpcClient() *Client {
	return &Client{
		rpcAddr: cc.Config().Heartbeat.Addr,
		Timeout: time.Duration(cc.Config().Heartbeat.Timeout) * time.Second,
	}
}

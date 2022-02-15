package connpool

import (
	"fmt"
	msgpackrpc "github.com/hashicorp/net-rpc-msgpackrpc"
	"net"
	"sync"
	"time"

	connp "github.com/toolkits/conn_pool"
	rpcpool "github.com/toolkits/conn_pool/rpc_conn_pool"
)

// SafeRpcConnPools ConnPools Manager 多个连接池管理
type SafeRpcConnPools struct {
	sync.RWMutex
	M           map[string]*connp.ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
}

func CreateSafeRpcConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRpcConnPools {
	cp := &SafeRpcConnPools{M: make(map[string]*connp.ConnPool), maxConns: maxConns, maxIdle: maxIdle,
		connTimeout: connTimeout, callTimeout: callTimeout}

	ct := time.Duration(cp.connTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneRpcPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

func createOneRpcPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
	p := connp.NewConnPool(name, address, int32(maxConns), int32(maxIdle))
	p.New = func(connName string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			return nil, err
		}

		return rpcpool.NewRpcClient(msgpackrpc.NewClient(conn), connName), nil
	}

	return p
}

// Call 同步发送, 完成发送或超时后 才能返回
func (s *SafeRpcConnPools) Call(addr, method string, args interface{}, resp interface{}) error {
	connPool, exists := s.Get(addr)
	if !exists {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err := connPool.Fetch()
	if err != nil {
		return fmt.Errorf("%s get connection fail: conn %v, err %v. proc: %s", addr, conn, err, connPool.Proc())
	}

	rpcClient := conn.(*rpcpool.RpcClient)
	callTimeout := time.Duration(s.callTimeout) * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(callTimeout):
		connPool.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", addr)
	case err = <-done:
		if err != nil {
			connPool.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", addr, err, connPool.Proc())
		} else {
			connPool.Release(conn)
		}
		return err
	}
}

func (s *SafeRpcConnPools) Get(address string) (*connp.ConnPool, bool) {
	s.RLock()
	defer s.RUnlock()
	p, exists := s.M[address]
	return p, exists
}

func (s *SafeRpcConnPools) Destroy() {
	s.Lock()
	defer s.Unlock()
	addresses := make([]string, 0, len(s.M))
	for address := range s.M {
		addresses = append(addresses, address)
	}

	for _, address := range addresses {
		s.M[address].Destroy()
		delete(s.M, address)
	}
}

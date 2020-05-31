package pool

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	connp "github.com/toolkits/conn_pool"
	rpcpool "github.com/toolkits/conn_pool/rpc_conn_pool"
)

// SafeRPCConnPools ConnPools Manager
type SafeRPCConnPools struct {
	sync.RWMutex
	M           map[string]*connp.ConnPool
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

// CreateSafeRPCConnPools TODO:
func CreateSafeRPCConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRPCConnPools {
	cp := &SafeRPCConnPools{M: make(map[string]*connp.ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneRPCPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

// CreateSafeJSONRPCConnPools TODO:
func CreateSafeJSONRPCConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *SafeRPCConnPools {
	cp := &SafeRPCConnPools{M: make(map[string]*connp.ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneJSONRPCPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

// Call 同步发送, 完成发送或超时后 才能返回
func (s *SafeRPCConnPools) Call(addr, method string, args interface{}, resp interface{}) error {
	connPool, exists := s.Get(addr)
	if !exists {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err := connPool.Fetch()
	if err != nil {
		return fmt.Errorf("%s get connection fail: conn %v, err %v. proc: %s", addr, conn, err, connPool.Proc())
	}

	rpcClient := conn.(*rpcpool.RpcClient)
	callTimeout := time.Duration(s.CallTimeout) * time.Millisecond

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

// Get TODO:
func (s *SafeRPCConnPools) Get(address string) (*connp.ConnPool, bool) {
	s.RLock()
	defer s.RUnlock()
	p, exists := s.M[address]
	return p, exists
}

// Destroy TODO:
func (s *SafeRPCConnPools) Destroy() {
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

// Proc TODO:
func (s *SafeRPCConnPools) Proc() []string {
	procs := []string{}
	for _, cp := range s.M {
		procs = append(procs, cp.Proc())
	}
	return procs
}

func createOneRPCPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
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

		return rpcpool.NewRpcClient(rpc.NewClient(conn), connName), nil
	}

	return p
}

func createOneJSONRPCPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
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

		return rpcpool.NewRpcClientWithCodec(jsonrpc.NewClientCodec(conn), connName), nil
	}

	return p
}

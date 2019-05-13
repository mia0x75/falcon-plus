package g

import (
	"errors"
	"math"
	"net/rpc"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net"
)

// SingleConnRPCClient TODO:
type SingleConnRPCClient struct {
	sync.Mutex
	rpcClient  *rpc.Client
	RPCServers []string
	Timeout    time.Duration
}

func (rpc *SingleConnRPCClient) close() {
	if rpc.rpcClient != nil {
		rpc.rpcClient.Close()
		rpc.rpcClient = nil
	}
}

func (rpc *SingleConnRPCClient) serverConn() error {
	if rpc.rpcClient != nil {
		return nil
	}

	var err error
	var retry int

	for _, addr := range rpc.RPCServers {
		retry = 1

	RETRY:
		rpc.rpcClient, err = net.JsonRpcClient("tcp", addr, rpc.Timeout)
		if err != nil {
			log.Errorf("[E] net.JsonRpcClient failed: %v", err)
			if retry > 3 {
				continue
			}

			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			goto RETRY
		}
		log.Infof("[I] connected RPC server: %s", addr)

		return nil
	}

	return errors.New("connect to RPC servers failed")
}

// Call TODO:
func (rpc *SingleConnRPCClient) Call(method string, args interface{}, reply interface{}) error {
	rpc.Lock()
	defer rpc.Unlock()

	err := rpc.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := rpc.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Warnf("[W] rpc call timeout %v => %v", rpc.rpcClient, rpc.RPCServers)
		rpc.close()
		return errors.New("rpc call timeout")
	case err := <-done:
		if err != nil {
			rpc.close()
			return err
		}
	}

	return nil
}

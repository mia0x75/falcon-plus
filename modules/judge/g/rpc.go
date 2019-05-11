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
	rpcClient   *rpc.Client
	RPCServers  []string
	Timeout     time.Duration
	CallTimeout time.Duration
}

func (rpc *SingleConnRPCClient) close() {
	if rpc.rpcClient != nil {
		rpc.rpcClient.Close()
		rpc.rpcClient = nil
	}
}

func (rpc *SingleConnRPCClient) insureConn() {
	if rpc.rpcClient != nil {
		return
	}

	var err error
	retry := 1

	for {
		if rpc.rpcClient != nil {
			return
		}

		for _, s := range rpc.RPCServers {
			rpc.rpcClient, err = net.JsonRpcClient("tcp", s, rpc.Timeout)
			if err == nil {
				return
			}

			log.Errorf("[E] dial %s fail: %s", s, err)
		}

		if retry > 6 {
			retry = 1
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
}

// Call TODO:
func (rpc *SingleConnRPCClient) Call(method string, args interface{}, reply interface{}) error {
	rpc.Lock()
	defer rpc.Unlock()

	rpc.insureConn()

	done := make(chan error, 1)
	go func() {
		done <- rpc.rpcClient.Call(method, args, reply)
	}()

	var err error

	select {
	case <-time.After(rpc.CallTimeout):
		err = errors.New("call hbs timeout")
	case err = <-done:
	}

	if err != nil {
		rpc.close()
	}

	return err
}

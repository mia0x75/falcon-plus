package g

import (
	"errors"
	"github.com/toolkits/net"
	"log"
	"math"
	"net/rpc"
	"sync"
	"time"
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient  *rpc.Client
	RpcServers []string
	Timeout    time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) serverConn() error {
	if this.rpcClient != nil {
		return nil
	}

	var err error
	var retry int

	for _, addr := range this.RpcServers {
		retry = 1

	RETRY:
		this.rpcClient, err = net.JsonRpcClient("tcp", addr, this.Timeout)
		if err != nil {
			log.Println("net.JsonRpcClient failed", err)
			if retry > 3 {
				continue
			}

			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			goto RETRY
		}
		log.Println("connected RPC server", addr)

		return nil
	}

	return errors.New("connect to RPC servers failed")
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	err := this.serverConn()
	if err != nil {
		return err
	}

	timeout := time.Duration(10 * time.Second)
	done := make(chan error, 1)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", this.rpcClient, this.RpcServers)
		this.close()
		return errors.New("rpc call timeout")
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}

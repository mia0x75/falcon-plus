package g

import (
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

var (
	TransferClientsLock *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients     map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func SendMetrics(metrics []*cmodel.MetricValue, resp *cmodel.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Config().Transfer.Addrs)) {
		addr := Config().Transfer.Addrs[i]

		c := getTransferClient(addr)
		if c == nil {
			c = initTransferClient(addr)
		}

		if updateMetrics(c, metrics, resp) {
			break
		}
	}
}

func initTransferClient(addr string) *SingleConnRpcClient {
	addrs := []string{
		addr,
	}
	var c *SingleConnRpcClient = &SingleConnRpcClient{
		RpcServers: addrs,
		Timeout:    time.Duration(Config().Transfer.Timeout) * time.Millisecond,
	}
	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	TransferClients[addr] = c

	return c
}

func updateMetrics(c *SingleConnRpcClient, metrics []*cmodel.MetricValue, resp *cmodel.TransferResponse) bool {
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Errorf("[E] call Transfer.Update fail: %v, metrics: %v", err, metrics)
		return false
	}
	return true
}

func getTransferClient(addr string) *SingleConnRpcClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()

	if c, ok := TransferClients[addr]; ok {
		return c
	}
	return nil
}

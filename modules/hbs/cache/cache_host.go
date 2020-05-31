package cache

import (
	"sync"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

// 每次心跳的时候agent把hostname汇报上来，经常要知道这个机器的hostid，把此信息缓存
// key: hostname value: hostid
type SafeHostMap struct {
	sync.RWMutex
	M map[string]int
}

var HostMap = &SafeHostMap{M: make(map[string]int)}

func (m *SafeHostMap) GetID(hostname string) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	id, exists := m.M[hostname]
	return id, exists
}

func (m *SafeHostMap) Init() {
	hosts, err := db.QueryHosts()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = hosts
}

type SafeMonitoredHosts struct {
	sync.RWMutex
	M map[int]*cm.Host
}

var MonitoredHosts = &SafeMonitoredHosts{M: make(map[int]*cm.Host)}

func (m *SafeMonitoredHosts) Get() map[int]*cm.Host {
	m.RLock()
	defer m.RUnlock()
	return m.M
}

func (m *SafeMonitoredHosts) Init() {
	hosts, err := db.QueryMonitoredHosts()
	if err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.M = hosts
}

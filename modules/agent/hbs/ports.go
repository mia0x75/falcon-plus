package hbs

import (
	"sync"
)

var (
	ports     []int64
	portsLock = new(sync.RWMutex)
)

func ReportPorts() []int64 {
	portsLock.RLock()
	defer portsLock.RUnlock()
	return ports
}

func CacheReportPorts(value []int64) {
	portsLock.Lock()
	defer portsLock.Unlock()
	ports = value
}

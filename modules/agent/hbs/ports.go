package hbs

import (
	"sync"
)

var (
	ports     []int64
	portsLock = new(sync.RWMutex)
)

// ReportPorts 获取暂存的本地端口数据
func ReportPorts() []int64 {
	portsLock.RLock()
	defer portsLock.RUnlock()
	return ports
}

// CacheReportPorts 把本地端口数据存储到变量中
func CacheReportPorts(value []int64) {
	portsLock.Lock()
	defer portsLock.Unlock()
	ports = value
}

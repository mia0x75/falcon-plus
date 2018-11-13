package hbs

import (
	"sync"
)

var (
	paths     []string
	pathsLock = new(sync.RWMutex)
)

func ReportPaths() []string {
	pathsLock.RLock()
	defer pathsLock.RUnlock()
	return paths
}

func CacheReportPaths(value []string) {
	pathsLock.Lock()
	defer pathsLock.Unlock()
	paths = value
}

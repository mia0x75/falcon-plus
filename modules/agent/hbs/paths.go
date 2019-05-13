package hbs

import (
	"sync"
)

var (
	paths     []string
	pathsLock = new(sync.RWMutex)
)

// ReportPaths TODO:
func ReportPaths() []string {
	pathsLock.RLock()
	defer pathsLock.RUnlock()
	return paths
}

// CacheReportPaths TODO:
func CacheReportPaths(value []string) {
	pathsLock.Lock()
	defer pathsLock.Unlock()
	paths = value
}

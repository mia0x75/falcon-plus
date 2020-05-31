package hbs

import (
	"sync"
)

var (
	duPaths     []string
	duPathsLock = new(sync.RWMutex)
)

func ReportDu() []string {
	duPathsLock.RLock()
	defer duPathsLock.RUnlock()
	return duPaths
}

func CacheReportDu(paths []string) {
	duPathsLock.Lock()
	defer duPathsLock.Unlock()
	duPaths = paths
}

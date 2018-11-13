package hbs

import (
	"sync"
)

var (
	sources     map[string]string
	sourcesLock = new(sync.RWMutex)
)

func ReportSources() map[string]string {
	sourcesLock.RLock()
	defer sourcesLock.RUnlock()
	return sources
}

func CacheReportSources(value map[string]string) {
	sourcesLock.Lock()
	defer sourcesLock.Unlock()
	sources = value
}
